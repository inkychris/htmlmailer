package main

import (
	"errors"
	"fmt"
	jsonyaml "github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/net/publicsuffix"
	"gopkg.in/mail.v2"
	"gopkg.in/robfig/cron.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func htmlmailer(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Usage: htmlmailer <config_file>")
	}
	yamlConfig, err := ioutil.ReadFile(args[0])
	if err != nil {
		log.Fatal(err)
	}
	config, err := ConfigFromBytes(yamlConfig)
	if err != nil {
		log.Fatal(err)
	}

	if config.Schedule == "" {
		err := config.Run()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	scheduler := cron.New()
	_, err = scheduler.AddFunc(config.Schedule, func(){
		err := config.Run()
		if err != nil {
			log.Print(err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	scheduler.Start()
	select{}
}

func main() {
	RootCmd := &cobra.Command{
		Use:   "htmlmailer <config_file>",
		Short: "Send the HTML response of a URL via email",
		Run: htmlmailer,
	}
	err := RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	Login struct {
		Credentials struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"credentials"`
		Form struct {
			Action string `yaml:"action"`
			UsernameField string `yaml:"username_field"`
			PasswordField string `yaml:"password_field"`
		} `yaml:"form"`
	} `yaml:"login"`
	TargetUrl string `yaml:"target_url"`
	Schedule string `yaml:"schedule"`
	Email struct {
		To []string `yaml:"to"`
		From string `yaml:"from"`
		Subject string `yaml:"subject"`
	} `yaml:"email"`
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
	if config.Email.From == "" {
		config.Email.From = config.SMTP.Username
	}
	return
}

func (config *Config) NewClient() (*http.Client, error) {
	cookiejarOptions := cookiejar.Options{PublicSuffixList: publicsuffix.List}
	jar, err := cookiejar.New(&cookiejarOptions)
	if err != nil {
		return nil, err
	}
	return &http.Client{Jar: jar}, nil
}

func (config *Config) NewEmailDialer() *mail.Dialer {
	dialer := mail.NewDialer(
		config.SMTP.Host,
		config.SMTP.Port,
		config.SMTP.Username,
		config.SMTP.Password,
	)
	dialer.StartTLSPolicy = mail.MandatoryStartTLS
	return dialer
}

func (config *Config) Run() error {
	client, err := config.NewClient()
	if err != nil {
		return err
	}

	if config.Login.Form.Action != "" {
		resp, err := client.PostForm(
			config.Login.Form.Action,
			url.Values{
				config.Login.Form.UsernameField: {config.Login.Credentials.Username},
				config.Login.Form.PasswordField: {config.Login.Credentials.Password},
			},
		)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 400 {
			return errors.New(fmt.Sprintf("Unsuccessful status code: %s (%s)", resp.Status, config.Login.Form.Action))
		}
	}

	resp, err := client.Get(config.TargetUrl)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New("Unsuccessful status code: " + resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(body) == 0 {
		return errors.New("Body of response is empty")
	}

	dialer := config.NewEmailDialer()
	sender, err := dialer.Dial()
	if err != nil {
		return err
	}

	for _, address := range config.Email.To {
		message := mail.NewMessage()
		message.SetHeader("From", config.Email.From)
		message.SetHeader("To", address)
		message.SetHeader("Subject", config.Email.Subject)
		message.SetBody("text/html", string(body))
		err = mail.Send(sender, message)
		if err != nil {
			log.Printf("Failed to send email to %s: %s", address, err.Error())
		}
	}
	return nil
}

const configSchema = `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "id": "html_mailer.json",
  "title": "HTML Mailer configuration schema",
  "description": "Configuration for HTML Mailer",
  "type": "object",
  "required": [
    "target_url",
    "email",
    "smtp"
  ],
  "additionalProperties": false,
  "properties": {
    "login": {
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "credentials": {
          "type": "object",
          "additionalProperties": false,
          "properties": {
            "username": {"type": "string"},
            "password": {"type": "string"}
          }
        },
        "form": {
          "type": "object",
          "additionalProperties": false,
          "properties": {
            "action": {"type": "string"},
            "username_field": {"type": "string"},
            "password_field": {"type": "string"}
          }
        }
      }
    },
    "target_url": {"type": "string"},
    "schedule": {"type": "string"},
    "email": {
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "to": {
          "type": "array",
          "items": {"$ref": "#/definitions/email_address"},
          "minItems": 1
        },
        "from": {"$ref": "#/definitions/email_address"},
        "subject": {"type": "string"}
      }
    },
    "smtp": {
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "host": {"type": "string"},
        "port": {"type": "integer"},
        "username": {"type": "string"},
        "password": {"type": "string"}
      }
    }
  },

  "definitions": {
    "email_address": {
      "type": "string",
      "pattern": "^(?i)[A-Z0-9._%+-]+@[A-Z0-9.-]+\\.[A-Z]{2,}$"
    }
  }
}
`
