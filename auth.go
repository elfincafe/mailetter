package mailetter

type Auth struct {
	user     string
	password string
}

func NewAuth(user string, password string) *Auth {
	a := new(Auth)
	a.user = user
	a.password = password
	return a
}
