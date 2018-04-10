package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"log"
	"os"
)

// For lazy initialization
var services *awsServices

type awsServices struct {
	ssm *ssm.SSM
}

func newAWSServices() *awsServices {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Fatalf("failed to create a new session.\n %v", err)
	}
	if aws.StringValue(sess.Config.Region) == "" {
		trace("no explicit region configurations. So now retrieving ec2metadata...")
		region, err := ec2metadata.New(sess).Region()
		if err != nil {
			trace(err)
			log.Fatalf("could not find region configurations")
		}
		sess.Config.Region = aws.String(region)
	}

	if arn := os.Getenv("ENV_INJECTOR_ASSUME_ROLE_ARN"); arn != "" {
		creds := stscreds.NewCredentials(sess, arn)
		return &awsServices{ssm: ssm.New(sess, &aws.Config{Credentials: creds})}
	} else {
		return &awsServices{ssm: ssm.New(sess)}
	}
}

// Initialize services lazily, and return it.
func getService() *awsServices {
	if services == nil {
		services = newAWSServices()
	}
	return services
}
