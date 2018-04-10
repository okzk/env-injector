package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("[env-injector] ")

	injectEnviron()

	args := os.Args
	if len(args) <= 1 {
		log.Fatal("missing command")
	}

	path, err := exec.LookPath(args[1])
	if err != nil {
		log.Fatal(err)
	}
	err = syscall.Exec(path, args[1:], os.Environ())
	if err != nil {
		log.Fatal(err)
	}
}

func injectEnviron() {
	injectEnvironByPath()
	injectEnvironByPrefix()
}
