package crawler

import (
	"net/http"
	"net/url"
)

type Client interface {
	Get(url string) (resp *http.Response, err error)
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}

type AuthClient interface {
	Get(url string) (resp *http.Response, err error)
	PostForm(url string, data url.Values) (resp *http.Response, err error)
	Login(credentials Credentials) (resp *http.Response, err error)
}

type LoginForm struct {
	Action        string
	UsernameField string
	PasswordField string
}

type session struct {
	client Client
	loginForm LoginForm
}

func Session(client Client, form LoginForm) *session {
	return &session{
		client: client,
		loginForm: form,
	}
}

func (session *session) Get(url string) (*http.Response, error) {
	return session.client.Get(url)
}

func (session *session) PostForm(url string, data url.Values) (*http.Response, error) {
	return session.client.PostForm(url, data)
}

func (session *session) Login(credentials Credentials) (resp *http.Response, err error) {
	return session.PostForm(
		session.loginForm.Action,
		url.Values{
			session.loginForm.UsernameField: {credentials.Username()},
			session.loginForm.PasswordField: {credentials.Password()},
		},
	)
}
