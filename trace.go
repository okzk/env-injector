package main

import (
	"log"
	"os"
)

var verbose = os.Getenv("ENV_INJECTOR_VERBOSE") == "1"

func trace(v ...interface{}) {
	if verbose {
		log.Println(v...)
	}
}

func tracef(format string, v ...interface{}) {
	if verbose {
		log.Printf(format, v...)
	}
}
