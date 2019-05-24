package crawler_test

import (
	"github.com/stretchr/testify/assert"
	"htmlmailer/crawler"
	"testing"
)

func TestNewCredentials(t *testing.T) {
	username := "*username*"
	password := "Password123"
	credentials := crawler.NewCredentials(username, password)
	assert.Equal(t, username, credentials.Username())
	assert.Equal(t, password, credentials.Password())
}
