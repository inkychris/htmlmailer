package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const CompleteConfig = `login:
  credentials:
    username: bob@example.com
    password: wxyz6789
  form:
    action: https://some.site.com/login
    username_field: user[email]
    password_field: user[password]
target_url: https://some.site.com/user/bob/example
schedule: 0 7 * * 2
email:
  to:
    - alice@example.com
    - charlie@example.com
  from: bob@example.com
  subject: Bob's Example

smtp:
  host: example.com
  port: 465
  username: bob@example.com
  password: abcd1234
`

const MinimalConfig = `
target_url: https://some.site.com/user/bob/example
email:
  to:
    - alice@example.com
smtp:
  host: example.com
  port: 465
  username: bob@example.com
  password: abcd1234
`

const ConfigMissingTarget = `
email:
  to:
    - alice@example.com
smtp:
  host: example.com
  port: 465
  username: bob@example.com
  password: abcd1234
`

const ConfigMissingEmail = `
target_url: https://some.site.com/user/bob/example
smtp:
  host: example.com
  port: 465
  username: bob@example.com
  password: abcd1234
`

const ConfigMissingSMTP = `
target_url: https://some.site.com/user/bob/example
email:
  to:
    - alice@example.com
`

const ConfigMissingEmailRecipient = `
target_url: https://some.site.com/user/bob/example
email:
  to: []
smtp:
  host: example.com
  port: 465
  username: bob@example.com
  password: abcd1234
`

func TestValidateConfigSchema_Complete(t *testing.T) {
	err := ValidateConfigSchema([]byte(CompleteConfig))
	require.NoError(t, err)
}

func TestValidateConfigSchema_Minimal(t *testing.T) {
	err := ValidateConfigSchema([]byte(MinimalConfig))
	require.NoError(t, err)
}

func TestValidateConfigSchema_MissingTarget(t *testing.T) {
	err := ValidateConfigSchema([]byte(ConfigMissingTarget))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "target_url")
}

func TestValidateConfigSchema_MissingEmail(t *testing.T) {
	err := ValidateConfigSchema([]byte(ConfigMissingEmail))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "email")
}

func TestValidateConfigSchema_MissingSMTP(t *testing.T) {
	err := ValidateConfigSchema([]byte(ConfigMissingSMTP))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "smtp")
}

func TestValidateConfigSchema_MissingEmailRecipient(t *testing.T) {
	err := ValidateConfigSchema([]byte(ConfigMissingEmailRecipient))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "email.to")
}

func TestConfigFromBytes_Complete(t *testing.T) {
	config, err := ConfigFromBytes([]byte(CompleteConfig))
	require.NoError(t, err)
	assert.Equal(t, "bob@example.com", config.Login.Credentials.Username)
	assert.Equal(t, "wxyz6789", config.Login.Credentials.Password)
	assert.Equal(t, "https://some.site.com/login", config.Login.Form.Action)
	assert.Equal(t, "user[email]", config.Login.Form.UsernameField)
	assert.Equal(t, "user[password]", config.Login.Form.PasswordField)
	assert.Equal(t, "https://some.site.com/user/bob/example", config.TargetUrl)
	assert.Equal(t, "0 7 * * 2", config.Schedule)
	assert.Equal(t, []string{"alice@example.com", "charlie@example.com"}, config.Email.To)
	assert.Equal(t, "bob@example.com", config.Email.From)
	assert.Equal(t, "Bob's Example", config.Email.Subject)
	assert.Equal(t, "example.com", config.SMTP.Host)
	assert.Equal(t, 465, config.SMTP.Port)
	assert.Equal(t, "bob@example.com", config.SMTP.Username)
	assert.Equal(t, "abcd1234", config.SMTP.Password)
}

func TestConfigFromBytes_Minimal(t *testing.T) {
	config, err := ConfigFromBytes([]byte(MinimalConfig))
	require.NoError(t, err)
	assert.Equal(t, config.Login, Config{}.Login)
	assert.Equal(t, "https://some.site.com/user/bob/example", config.TargetUrl)
	assert.Empty(t, config.Schedule)
	assert.Equal(t, []string{"alice@example.com"}, config.Email.To)
	assert.Equal(t, "bob@example.com", config.Email.From)
	assert.Empty(t, config.Email.Subject)
	assert.Equal(t, "example.com", config.SMTP.Host)
	assert.Equal(t, 465, config.SMTP.Port)
	assert.Equal(t, "bob@example.com", config.SMTP.Username)
	assert.Equal(t, "abcd1234", config.SMTP.Password)
}

func TestConfigFromBytes_MissingTarget(t *testing.T) {
	_, err := ConfigFromBytes([]byte(ConfigMissingTarget))
	require.Error(t, err)
}
