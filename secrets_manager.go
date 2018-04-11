package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"log"
	"os"
)

func injectEnvironSecretManager(name string) {
	tracef("secret name: %s", name)

	svc := getService().secretsManager
	ret, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(name),
	})
	if err != nil {
		log.Fatalf("secretsmanager:GetSecretValue failed. (name: %s)\n %v", name, err)
	}
	secrets := make(map[string]string)
	err = json.Unmarshal([]byte(aws.StringValue(ret.SecretString)), &secrets)
	if err != nil {
		log.Fatalf("secretsmanager:GetSecretValue returns invalid json. (name: %s)\n %v", name, err)
	}
	for key, val := range secrets {
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
			tracef("env injected: %s", key)
		}
	}
}
