package mailetter

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Dsn struct {
	scheme string
	host   string
	port   int
}

func NewDsn(str string) (*Dsn, error) {
	str = strings.ToLower(strings.TrimSpace(str))
	u, err := url.Parse(str)
	if err != nil {
		return nil, err
	}
	dsn := new(Dsn)
	defaultPort := 0
	switch u.Scheme {
	case "smtp":
		dsn.scheme = "smtp"
		defaultPort = 25
	case "smtps":
		dsn.scheme = "smtps"
		defaultPort = 465
	case "smtp+tls":
		dsn.scheme = "smtp+tls"
		defaultPort = 587
	default:
		dsn.scheme = "smtps"
		defaultPort = 465
	}
	dsn.host = u.Hostname()
	if dsn.host == "" {
		return nil, fmt.Errorf(`Empty Hostname`)
	}
	if u.Port() != "" {
		port, err := strconv.ParseInt(u.Port(), 10, 32)
		if err != nil {
			return nil, err
		}
		dsn.port = int(port)
	} else {
		dsn.port = defaultPort
	}
	return dsn, nil
}

func (dsn *Dsn) Scheme() string {
	return dsn.scheme
}

func (dsn *Dsn) Host() string {
	return dsn.host
}

func (dsn *Dsn) Port() int {
	return dsn.port
}

func (dsn *Dsn) Socket() string {
	return fmt.Sprintf("%s:%d", dsn.host, dsn.port)
}
