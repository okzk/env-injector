package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var verbose = os.Getenv("ENV_INJECTOR_VERBOSE") == "1"

func main() {
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

func injectEnviron() {
	prefix := os.Getenv("ENV_INJECTOR_PREFIX")
	if prefix == "" {
		trace("no parameter prefix specified, skipping injection")
		return
	}
	tracef("parameter prefix: %s", prefix)
	if !strings.HasSuffix(prefix, ".") {
		prefix += "."
	}

	names := []*string{}
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 && parts[1] == "" {
			names = append(names, aws.String(prefix+parts[0]))
		}
	}
	if len(names) == 0 {
		trace("nothing to be injected")
		return
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		trace(err)
		trace("failed to create session")
		return
	}
	if *sess.Config.Region == "" {
		trace("no explict region configuration. So now retriving ec2metadata...")
		region, err := ec2metadata.New(sess).Region()
		if err != nil {
			trace(err)
			trace("could not find region configuration")
			return
		}
		sess.Config.Region = aws.String(region)
	}

	svc := ssm.New(sess)
	result, err := svc.GetParameters(&ssm.GetParametersInput{
		Names:          names,
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		trace(err)
		return
	}

	if len(result.InvalidParameters) != 0 {
		tracef("invalid parameters: %v", result.InvalidParameters)
	}
	for _, param := range result.Parameters {
		key := strings.TrimPrefix(*param.Name, prefix)
		os.Setenv(key, *param.Value)
		tracef("env injected: %s", key)
	}
}
