package main

import (
	"os"

	ps "github.com/mitchellh/go-ps"
	log "github.com/sirupsen/logrus"
)

// CommitString is the commit used to build the server
var CommitString string

var except = map[string]struct{}{
	"remote":    {},
	"syncthing": {},
}

func shouldKill(p ps.Process) bool {
	if _, ok := except[p.Executable()]; ok {
		return false
	}

	// don't kill the root process of the container
	if p.PPid() == 0 {
		return false
	}

	// don't kill your parent process
	if p.Pid() == os.Getppid() {
		return false
	}

	// don't kill yourself
	if p.Pid() == os.Getpid() {
		return false
	}

	return true
}

func main() {
	log.Infof("clean service %s started", CommitString)
	processes, err := ps.Processes()
	if err != nil {
		log.Errorf("fail to list processes: %s", processes)
		os.Exit(1)
	}

	for _, p := range processes {

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
