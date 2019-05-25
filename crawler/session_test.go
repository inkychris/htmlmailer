package crawler_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"htmlmailer/crawler"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

type FakeClient struct {
	NoAuthGetUrl   string
	AuthGetUrl     string
	PostUrl        string
	PostFormValues url.Values
	Credentials    crawler.Credentials
	LoginForm      crawler.LoginForm
}

func (client *FakeClient) IsAuthenticated() bool {
	if client.PostFormValues == nil {
		return false
	}
	if ! (client.PostFormValues[client.LoginForm.UsernameField][0] == client.Credentials.Username()) {
		return false
	}
	return client.PostFormValues[client.LoginForm.PasswordField][0] == client.Credentials.Password()
}

func (client *FakeClient) Get(url string) (resp *http.Response, err error) {
	switch url {
	case client.NoAuthGetUrl:
		return &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("Success"))),
		}, nil
	case client.AuthGetUrl:

		return &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("Success"))),
		}, nil
	default:
		return &http.Response{
			Status:     "404 PAGE NOT FOUND",
			StatusCode: 404,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("Page not found"))),
		}, nil
	}
}

func (client *FakeClient) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	client.PostFormValues = data
	if url == client.PostUrl {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("Success"))),
		}, nil
	}
	return &http.Response{
		StatusCode: 404,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte("Page not found"))),
	}, nil
}

func NewTestSession() (session *crawler.Session, client *FakeClient, credentials crawler.Credentials) {
	credentials = crawler.NewCredentials("testuser", "Password123")
	loginForm := crawler.LoginForm{
		Url: "example.com/login",
		UsernameField: "login[user]",
		PasswordField: "login[pass]",
	}
	client = &FakeClient{
		NoAuthGetUrl: "example.com/get",
		AuthGetUrl:   "example.com/user/get",
		PostUrl:      "example.com/login",
		LoginForm:    loginForm,
		Credentials:  credentials,
	}
	session = crawler.NewSession(client, loginForm)
	return
}

func TestSession_NoAuthGet(t *testing.T) {
	session, client, _ := NewTestSession()
	assert.False(t, client.IsAuthenticated())
	resp, err := session.Get(client.NoAuthGetUrl)
	require.NoError(t, err)
	assert.Equal(t,200, resp.StatusCode)
}

func TestSession_Login(t *testing.T) {
	session, client, credentials := NewTestSession()
	assert.False(t, client.IsAuthenticated())
	resp, err := session.Login(credentials)
	require.NoError(t, err)
	assert.Equal(t,200, resp.StatusCode)
	assert.True(t, client.IsAuthenticated())
}

func TestSession_AuthGet(t *testing.T) {
	session, client, credentials := NewTestSession()
	assert.False(t, client.IsAuthenticated())
	resp, err := session.Login(credentials)
	require.NoError(t, err)
	assert.Equal(t,200, resp.StatusCode)
	assert.True(t, client.IsAuthenticated())
	resp, err = session.Get(client.AuthGetUrl)
	require.NoError(t, err)
	assert.Equal(t,200, resp.StatusCode)
}
