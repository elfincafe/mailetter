package mailetter

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewDsn(t *testing.T) {
	type tcase struct {
		dsn    string
		typ    string
		errmsg string
	}
	cases := []tcase{}
	cases = append(cases, tcase{"smtp://localhost:25", "*mailetter.Dsn", ""})
	cases = append(cases, tcase{"smtps://localhost:465", "*mailetter.Dsn", ""})
	cases = append(cases, tcase{"smtp+tls://localhost:587", "*mailetter.Dsn", ""})
	cases = append(cases, tcase{"smtp://localhost", "*mailetter.Dsn", ""})
	cases = append(cases, tcase{"smtps://localhost", "*mailetter.Dsn", ""})
	cases = append(cases, tcase{"smtp+tls://localhost", "*mailetter.Dsn", ""})
	cases = append(cases, tcase{"aaa:/", "", "Empty Hostname"})
	cases = append(cases, tcase{"::::", "", "missing protocol scheme"})
	for k, v := range cases {
		dsn, err := NewDsn(v.dsn)
		if err != nil && !strings.Contains(err.Error(), v.errmsg) {
			t.Errorf("[Case%d] ErrMsg: %s (%s)", k, err.Error(), v.errmsg)
		}
		typ := reflect.TypeOf(dsn)
		if err == nil && typ.String() != v.typ {
			t.Errorf("[Case%d] Type: %s != %s", k, typ.String(), v.typ)
		}
	}
}

func TestScheme(t *testing.T) {
	type tcase struct {
		dsn    string
		scheme string
	}
	cases := []tcase{}
	cases = append(cases, tcase{"smtp://localhost:25", "smtp"})
	cases = append(cases, tcase{"smtps://localhost:465", "smtps"})
	cases = append(cases, tcase{"smtp+tls://localhost:587", "smtp+tls"})
	cases = append(cases, tcase{"SMTP://localhost:25", "smtp"})
	cases = append(cases, tcase{"SMTPS://localhost:465", "smtps"})
	cases = append(cases, tcase{"SMTP+TLS://localhost:587", "smtp+tls"})
	cases = append(cases, tcase{"enigma://localhost:587", "smtps"})
	cases = append(cases, tcase{"ENIGMA://localhost:587", "smtps"})
	for k, v := range cases {
		dsn, _ := NewDsn(v.dsn)
		if dsn.Scheme() != v.scheme {
			t.Errorf("[Case%d] %s (%s)", k, dsn.Scheme(), v.scheme)
		}
	}
}

func TestHost(t *testing.T) {
	type tcase struct {
		dsn  string
		host string
	}
	cases := []tcase{}
	cases = append(cases, tcase{"smtp://localhost:25", "localhost"})
	cases = append(cases, tcase{"smtps://example.com", "example.com"})
	for k, v := range cases {
		dsn, _ := NewDsn(v.dsn)
		if dsn.Host() != v.host {
			t.Errorf("[Case%d] %s (%s)", k, dsn.Host(), v.host)
		}
	}
}

func TestPort(t *testing.T) {
	type tcase struct {
		dsn  string
		port int
	}
	cases := []tcase{}
	cases = append(cases, tcase{"smtp://localhost:25", 25})
	cases = append(cases, tcase{"smtp://localhost:10025", 10025})
	cases = append(cases, tcase{"smtp://localhost", 25})
	cases = append(cases, tcase{"smtps://localhost:465", 465})
	cases = append(cases, tcase{"smtps://localhost:10465", 10465})
	cases = append(cases, tcase{"smtps://localhost", 465})
	cases = append(cases, tcase{"smtp+tls://localhost:587", 587})
	cases = append(cases, tcase{"smtp+tls://localhost:10587", 10587})
	cases = append(cases, tcase{"smtp+tls://localhost", 587})
	cases = append(cases, tcase{"enigma://localhost:465", 465})
	cases = append(cases, tcase{"enigma://localhost:10465", 10465})
	cases = append(cases, tcase{"enigma://localhost", 465})
	for k, v := range cases {
		dsn, _ := NewDsn(v.dsn)
		if dsn.Port() != v.port {
			t.Errorf("[Case%d] Port: %d (%d)", k, dsn.Port(), v.port)
		}
	}
}

func TestSocket(t *testing.T) {
	type tcase struct {
		dsn    string
		socket string
	}
	cases := []tcase{}
	cases = append(cases, tcase{"smtp://localhost:25", "localhost:25"})
	cases = append(cases, tcase{"smtp://localhost:10025", "localhost:10025"})
	cases = append(cases, tcase{"smtp://localhost", "localhost:25"})
	cases = append(cases, tcase{"smtps://localhost:465", "localhost:465"})
	cases = append(cases, tcase{"smtps://localhost:10465", "localhost:10465"})
	cases = append(cases, tcase{"smtps://localhost", "localhost:465"})
	cases = append(cases, tcase{"smtp+tls://localhost:587", "localhost:587"})
	cases = append(cases, tcase{"smtp+tls://localhost:10587", "localhost:10587"})
	cases = append(cases, tcase{"smtp+tls://localhost", "localhost:587"})
	cases = append(cases, tcase{"enigma://localhost:465", "localhost:465"})
	cases = append(cases, tcase{"enigma://localhost:10465", "localhost:10465"})
	cases = append(cases, tcase{"enigma://localhost", "localhost:465"})
	for k, v := range cases {
		dsn, _ := NewDsn(v.dsn)
		if dsn.Socket() != v.socket {
			t.Errorf("[Case%d] Socket: %s (%s)", k, dsn.Socket(), v.socket)
		}
	}
}
