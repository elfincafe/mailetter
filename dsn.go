package mailetter

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type dsn struct {
	scheme string
	host   string
	port   int
}

func newDsn(str string) (*dsn, error) {
	str = strings.ToLower(strings.TrimSpace(str))
	u, err := url.Parse(str)
	if err != nil {
		return nil, err
	}
	sDsn := new(dsn)
	defaultPort := 0
	switch u.Scheme {
	case "smtp":
		sDsn.scheme = "smtp"
		defaultPort = 25
	case "smtps":
		sDsn.scheme = "smtps"
		defaultPort = 465
	case "smtp+tls":
		sDsn.scheme = "smtp+tls"
		defaultPort = 587
	default:
		sDsn.scheme = "smtps"
		defaultPort = 465
	}
	sDsn.host = u.Hostname()
	if sDsn.host == "" {
		return nil, fmt.Errorf(`Empty Hostname`)
	}
	if u.Port() != "" {
		port, err := strconv.ParseInt(u.Port(), 10, 32)
		if err != nil {
			return nil, err
		}
		sDsn.port = int(port)
	} else {
		sDsn.port = defaultPort
	}
	return sDsn, nil
}

func (dsn *dsn) Scheme() string {
	return dsn.scheme
}

func (dsn *dsn) Host() string {
	return dsn.host
}

func (dsn *dsn) Port() int {
	return dsn.port
}

func (dsn *dsn) Socket() string {
	return fmt.Sprintf("%s:%d", dsn.host, dsn.port)
}

func (dsn *dsn) IsSsl() bool {
	return dsn.scheme == "smtps"
}
