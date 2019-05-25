package crawler

import (
	"net/http"
	"net/url"
)

type Client interface {
	Get(url string) (resp *http.Response, err error)
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}

type LoginForm struct {
	Url string
	UsernameField string
	PasswordField string
}

type Session struct {
	client Client
	loginForm LoginForm
}

func NewSession(client Client, form LoginForm) *Session {
	return &Session{
		client: client,
		loginForm: form,
	}
}

func (session *Session) Get(url string) (*http.Response, error) {
	return session.client.Get(url)
}

func (session *Session) PostForm(url string, data url.Values) (*http.Response, error) {
	return session.client.PostForm(url, data)
}

func (session *Session) Login(credentials Credentials) (resp *http.Response, err error) {
	return session.PostForm(
		session.loginForm.Url,
		url.Values{
			session.loginForm.UsernameField: {credentials.Username()},
			session.loginForm.PasswordField: {credentials.Password()},
		},
	)
}
