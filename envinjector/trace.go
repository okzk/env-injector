package envinjector

import (
	"log"
	"os"
)

var verbose = os.Getenv("ENV_INJECTOR_VERBOSE") == "1"
var logger = log.New(os.Stderr, "[env-injector] ", 0)

// ConfigLogger calls the given function with internal logger.
func ConfigLogger(f func(logger *log.Logger)) {
	f(logger)
}

func trace(v ...interface{}) {
	if verbose {
		logger.Println(v...)
	}
}

func tracef(format string, v ...interface{}) {
	if verbose {
		logger.Printf(format, v...)
	}
}
