package mailetter

import (
	"fmt"
	"net/url"
	"strings"
)

type Dsn struct {
	Scheme string
	Host   string
	Port   string
	User   *url.Userinfo
}

func NewDsn(dsnStr string) (*Dsn, error) {
	dsnStr = strings.ToLower(dsnStr)
	u, err := url.Parse(dsnStr)
	if err != nil {
		return nil, err
	}
	dsn := new(Dsn)
	switch u.Scheme {
	case "smtp":
		dsn.Scheme = "smtp"
	case "smtps":
		dsn.Scheme = "smtps"
	case "smtp+tls":
		dsn.Scheme = "smtp+tls"
	default:
		return nil, fmt.Errorf("DSN is invalid")
	}
	dsn.Host = u.Hostname()
	dsn.Port = u.Port()
	dsn.User = u.User
	return dsn, nil
}
