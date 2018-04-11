package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"gopkg.in/yaml.v2"
	"log"
	"strings"
)

type envKeyDecorator interface {
	decorate(key string) string
}

type noopKeyDecorator struct{}

func (*noopKeyDecorator) decorate(key string) string { return key }

var noop = &noopKeyDecorator{}

type metaConfig struct {
	SecretName         string `yaml:"secret_name"`
	ParameterStorePath string `yaml:"parameter_store_path"`

	EnvPrefix  string `yaml:"env_prefix"`
	Capitalize bool   `yaml:"capitalize"`
}

func (m *metaConfig) decorate(key string) string {
	if m.EnvPrefix != "" {
		key = m.EnvPrefix + "_" + key
	}
	if m.Capitalize {
		return strings.ToUpper(key)
	}

	return key
}

func (m *metaConfig) injectEnviron() {
	if m.SecretName != "" {
		injectEnvironSecretManager(m.SecretName, m)
	}
	if m.ParameterStorePath != "" {
		injectEnvironByPath(m.ParameterStorePath, m)
	}
}

func injectEnvironViaMetaConfig(path string) {
	svc := getService().ssm
	result, err := svc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(path),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		log.Fatalf("ssm:GetParameter failed. (path: %s)\n %v", path, err)
	}

	configs := make([]*metaConfig, 0)
	err = yaml.Unmarshal([]byte(aws.StringValue(result.Parameter.Value)), &configs)
	if err != nil {
		log.Fatalf("ssm:GetParameter returns invalid yaml. (path: %s)\n %v", path, err)
	}

	for _, config := range configs {
		config.injectEnviron()
	}
}
