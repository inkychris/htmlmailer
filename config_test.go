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
