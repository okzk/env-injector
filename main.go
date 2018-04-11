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
	if name := os.Getenv("ENV_INJECTOR_SECRET_NAME"); name != "" {
		injectEnvironSecretManager(name)
	} else {
		trace("no secret name specified, skipping injection by SecretsManager")
	}

	if path := os.Getenv("ENV_INJECTOR_PATH"); path != "" {
		injectEnvironByPath(path)
	} else {
		trace("no parameter path specified, skipping injection by path")
	}

	if prefix := os.Getenv("ENV_INJECTOR_PREFIX"); prefix != "" {
		injectEnvironByPrefix(prefix)
	} else {
		trace("no parameter prefix specified, skipping injection by prefix")
	}
}
