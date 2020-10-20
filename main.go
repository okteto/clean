package main

import (
	"os"

	ps "github.com/mitchellh/go-ps"
	log "github.com/sirupsen/logrus"
)

// CommitString is the commit used to build the server
var CommitString string

var except = map[string]struct{}{
	"okteto-remote": {},
	"syncthing":     {},
	"screen":        {},
	"tmux: server":  {},
}

func shouldKill(p ps.Process) bool {
	if p.Pid() == 1 {
		log.Info("not killing root process of the container")
		return false
	}

	if p.Pid() == os.Getppid() {
		log.Info("not killing parent process")
		return false
	}

	if p.Pid() == os.Getpid() {
		log.Info("not killing yourself")
		return false
	}

	if isChildrenOfExcept(p) {
		return false
	}

	return true
}

func isChildrenOfExcept(p ps.Process) bool {
	if p.Pid() == 1 {
		return false
	}

	if _, ok := except[p.Executable()]; ok {
		log.Infof("not killing, children of %s", p.Executable())
		return true
	}

	parent, err := ps.FindProcess(p.PPid())
	if err != nil {
		log.Errorf("fail to find process %d : %s", p.PPid(), err)
		return false
	}

	if parent == nil {
		return false
	}

	return isChildrenOfExcept(parent)
}

func main() {
	log.Infof("clean service started sha=%s pid=%d ppid=%d", CommitString, os.Getpid(), os.Getppid())
	processes, err := ps.Processes()
	if err != nil {
		log.Errorf("fail to list processes: %s", processes)
		os.Exit(1)
	}

	for _, p := range processes {
		log.Infof("checking process %s", p.Executable())

		if !shouldKill(p) {
			continue
		}

		pr, err := os.FindProcess(p.Pid())
		if err != nil {
			log.Errorf("fail to find process %d : %s", p.Pid(), err)
			continue
		}

		if err := pr.Kill(); err != nil {
			log.Errorf("fail to kill process %d : %s", p.Pid(), err)
		}
	}

	log.Infof("clean service finished")
}
