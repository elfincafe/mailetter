package mailetter

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type dsn struct {
	dsn    string
	scheme string
	host   string
	port   int
}

func newDsn(dsnStr string) *dsn {
	d := new(dsn)
	d.dsn = strings.ToLower(strings.TrimSpace(dsnStr))
	d.scheme = ""
	d.host = ""
	d.port = 0
	return d
}

func (d *dsn) parse() error {
	u, err := url.Parse(d.dsn)
	if err != nil {
		return err
	}
	defaultPort := 0
	switch u.Scheme {
	case "smtp":
		d.scheme = "smtp"
		defaultPort = 25
	case "smtps":
		d.scheme = "smtps"
		defaultPort = 465
	case "smtp+tls":
		d.scheme = "smtp+tls"
		defaultPort = 25
	}
	if d.scheme == "" {
		return fmt.Errorf(`the scheme must be one of "smtp", "smtps", or "smtp+tls"`)
	}
	d.host = u.Hostname()
	if d.host == "" {
		return fmt.Errorf(`the host must be ether hostname or ip address`)
	}
	if u.Port() != "" {
		port, err := strconv.ParseInt(u.Port(), 10, 32)
		if err != nil || port < 1 {
			return fmt.Errorf(`the port must be positive integer`)
		}
		d.port = int(port)
	} else {
		d.port = defaultPort
	}
	return nil
}

func (dsn *dsn) Socket() string {
	return fmt.Sprintf("%s:%d", dsn.host, dsn.port)
}

func (dsn *dsn) IsSsl() bool {
	return dsn.scheme == "smtps"
}
