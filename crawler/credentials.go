package crawler

type Credentials struct {
	username string
	password string
}

func NewCredentials(username string, password string) *Credentials {
	return &Credentials{
		username: username,
		password: password,
	}
}

func (credentials *Credentials) Username() string {
	return credentials.username
}

func (Credentials *Credentials) Password() string {
	return Credentials.password
}
