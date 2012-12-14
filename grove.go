package main

import (
	"io"
	"log"
	"net/http"
	"net/http/cgi"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	Version     = "0.2.1"
	DefaultPort = "8860"
)

const (
	usage          = "usage: %s [repositorydir]\n"
	gitHttpBackend = "git-http-backend"
)

var (
	g *GitBackendHandler
)

type GitBackendHandler struct {
	Handler *cgi.Handler
	Logger  *log.Logger
}

func main() {
	logger := log.New(os.Stdout, "", log.Ltime)

	var repodir string
	if len(os.Args) > 1 {
		repodir = os.Args[1]
		if !path.IsAbs(repodir) {
			wd, err := os.Getwd()
			if err != nil {
				logger.Fatalln("Error getting working directory:", err)
			}
			path.Join(wd, repodir)
		}
	} else {
		wd, err := os.Getwd()
		if err != nil {
			logger.Fatalln("Error getting working directory:", err)
		}
		repodir = wd
	}

	Serve(logger, repodir, DefaultPort)
}

func Serve(logger *log.Logger, repodir string, port string) (err error) {
	//Use 'git --exec-path' to get the path
	//of the git executables.
	var execPath []byte
	gitExecCmd := exec.Command("git", "--exec-path")
	execPath, err = gitExecCmd.Output()
	if err != nil {
		return
	}

	g = &GitBackendHandler{
		Handler: &cgi.Handler{
			Path:   strings.TrimRight(string(execPath), "\r\n") + "/" + gitHttpBackend,
			Root:   "/",
			Dir:    repodir,
			Env:    []string{"GIT_PROJECT_ROOT=" + repodir, "GIT_HTTP_EXPORT_ALL=TRUE"},
			Logger: logger,
		},
		Logger: logger,
	}
	logger.Println("Created CGI handler:",
		"\n\tPath:\t", g.Handler.Path,
		"\n\tRoot:\t", g.Handler.Root,
		"\n\tDir:\t", g.Handler.Dir,
		"\n\tEnv:\t",
		"\n\t\t", g.Handler.Env[0],
		"\n\t\t", g.Handler.Env[1])

	logger.Println("Starting server")
	http.HandleFunc("/", HandleWeb)
	err = http.ListenAndServe(":"+port, nil)
	return
}

func HandleWeb(w http.ResponseWriter, req *http.Request) {
	//Send the request to the git http backend
	//if it is to a .git URL.
	if strings.Contains(req.URL.String(), ".git") {
		g.Logger.Println("Git request to", req.URL, "from", req.RemoteAddr)
		g.Handler.ServeHTTP(w, req)
		return
	} else {
		g.Logger.Println("View of", req.URL, "from", req.RemoteAddr)
	}

	path := path.Join(g.Handler.Dir, req.URL.String())
	io.WriteString(w, ShowPath(path))
}