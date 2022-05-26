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
	port   uint16
}

func NewDsn(str string) (*Dsn, error) {
	str = strings.ToLower(strings.TrimSpace(str))
	u, err := url.Parse(str)
	if err != nil {
		return nil, err
	}
	dsn := new(Dsn)
	switch u.Scheme {
	case "smtp":
		dsn.scheme = "smtp"
	case "smtps":
		dsn.scheme = "smtps"
	case "smtp+tls":
		dsn.scheme = "smtp+tls"
	default:
		dsn.scheme = "smtps"
	}
	dsn.host = u.Hostname()
	println("----------------->")
	println(u.Hostname())
	println("<-----------------")
	if dsn.host == "" {
		return nil, fmt.Errorf(`Invalid Hostname %s`, dsn.host)
	}
	if u.Port() != "" {
		port, err := strconv.ParseUint(u.Port(), 10, 16)
		if err != nil {
			return nil, err
		}
		dsn.port = uint16(port)
	} else {
		dsn.port = uint16(587)
	}
	return dsn, nil
}

func (dsn *Dsn) Scheme() string {
	return dsn.scheme
}

func (dsn *Dsn) Host() string {
	return dsn.host
}

func (dsn *Dsn) Port() uint16 {
	return dsn.port
}

func (dsn *Dsn) String() string {
	return fmt.Sprintf("%s:%d", dsn.host, dsn.port)
}
