package main

import (
	"os/exec"
	"strings"
)

var (
	execPath string //The path containing git binaries
)

func setExecPath() (err error) {
	//Use 'git --exec-path' to get the path
	//of the git executables.
	path, err := execute("", "git", "--exec-path")
	execPath = strings.TrimRight(path, "\r\n")
	return
}

func gitBranch(path string) (branch string) {
	branch, _ = execute(path, "git", "rev-parse", "--abbrev-ref", "HEAD")
	return strings.TrimRight(branch, "\r\n")
}

func execute(dir, command string, args ...string) (output string, err error) {
	cmd := exec.Command(command, args...)
	if len(dir) != 0 {
		cmd.Dir = dir
	}
	out, err := cmd.Output()
	return string(out), err
}