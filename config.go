package main

import (
	"errors"
	"fmt"
	jsonyaml "github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
	"strings"
)

type Config struct {
	Login struct {
		Credentials struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"credentials"`
		Form struct {
			Host string `yaml:"host"`
			Port int `yaml:"port"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			From string `yaml:"from"`
		} `yaml:"form"`
	} `yaml:"login"`
	TargetUrl string `yaml:"target_url"`
	Schedule string `yaml:"schedule"`
	Email struct {} `yaml:"email"`
	SMTP struct {
		Host string `yaml:"host"`
		Port int `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"smtp"`
}

func ValidateConfigSchema(config []byte) (err error) {
	jsonConfig, err := jsonyaml.YAMLToJSON(config)
	if err != nil {
		return
	}
	schemaLoader := gojsonschema.NewStringLoader(configSchema)
	documentLoader := gojsonschema.NewBytesLoader(jsonConfig)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return
	}
	if !result.Valid() {
		schemaErrors := strings.Builder{}
		for _, schemaError := range result.Errors() {
			schemaErrors.WriteString(fmt.Sprintf("%s: %s\n", schemaError.Field(), schemaError.Description()))
		}
		return errors.New(schemaErrors.String())
	}
	return
}

func ConfigFromBytes(yamlConfig []byte) (config *Config, err error) {
	err = ValidateConfigSchema(yamlConfig)
	if err != nil {
		return nil, err
	}
	config = &Config{}
	err = yaml.Unmarshal(yamlConfig, config)
	return
}
