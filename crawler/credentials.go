package crawler

type Credentials interface {
	Username() string
	Password() string
}

type credentials struct {
	username string
	password string
}

func NewCredentials(username string, password string) Credentials {
	return &credentials{
		username: username,
		password: password,
	}
}

func (credentials *credentials) Username() string {
	return credentials.username
}

func (Credentials *credentials) Password() string {
	return Credentials.password
}
