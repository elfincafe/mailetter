package mailetter

import "net/smtp"

type (
	AuthInterface interface {
		smtp.Auth
	}
	PlainAuth struct {
		auth smtp.Auth
	}
	CramMd5Auth struct {
		auth smtp.Auth
	}
)
