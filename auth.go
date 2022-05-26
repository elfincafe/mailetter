package mailetter

type Auth struct {
	user     string
	password string
	valid    bool
}

func NewAuth(user string, password string) *Auth {
	a := new(Auth)
	a.user = user
	a.password = password
	if len(a.user) > 0 && len(a.password) > 0 {
		a.valid = true
	}
	return a
}
