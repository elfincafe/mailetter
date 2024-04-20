package mailetter

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewDsn(t *testing.T) {
	cases := []struct {
		dsn    string
		typ    string
		errmsg string
	}{
		{
			"smtp://localhost:25",
			"*mailetter.Dsn",
			"",
		},
		{
			"smtps://localhost:465",
			"*mailetter.Dsn",
			"",
		},
		{
			"smtp+tls://localhost:587",
			"*mailetter.Dsn",
			"",
		},
		{
			"smtp://localhost",
			"*mailetter.Dsn",
			"",
		},
		{
			"smtps://localhost",
			"*mailetter.Dsn",
			"",
		},
		{
			"smtp+tls://localhost",
			"*mailetter.Dsn",
			"",
		},
		{
			"aaa:/",
			"",
			"Empty Hostname",
		},
		{
			"::::",
			"",
			"missing protocol scheme",
		},
	}

	for k, v := range cases {
		dsn, err := newDsn(v.dsn)
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
	cases := []struct {
		dsn    string
		scheme string
	}{
		{
			"smtp://localhost:25",
			"smtp",
		},
		{
			"smtps://localhost:465",
			"smtps",
		},
		{
			"smtp+tls://localhost:587",
			"smtp+tls",
		},
		{
			"SMTP://localhost:25",
			"smtp",
		},
		{
			"SMTPS://localhost:465",
			"smtps",
		},
		{
			"SMTP+TLS://localhost:587",
			"smtp+tls",
		},
		{
			"enigma://localhost:587",
			"smtps",
		},
		{
			"ENIGMA://localhost:587",
			"smtps",
		},
	}

	for k, v := range cases {
		dsn, _ := newDsn(v.dsn)
		if dsn.Scheme() != v.scheme {
			t.Errorf("[Case%d] %s (%s)", k, dsn.Scheme(), v.scheme)
		}
	}
}

func TestHost(t *testing.T) {
	cases := []struct {
		dsn  string
		host string
	}{
		{
			"smtp://localhost:25",
			"localhost",
		},
		{
			"smtps://example.com",
			"example.com",
		},
	}

	for k, v := range cases {
		dsn, _ := newDsn(v.dsn)
		if dsn.Host() != v.host {
			t.Errorf("[Case%d] %s (%s)", k, dsn.Host(), v.host)
		}
	}
}

func TestPort(t *testing.T) {
	cases := []struct {
		dsn  string
		port int
	}{
		{
			"smtp://localhost:25",
			25,
		},
		{
			"smtp://localhost:10025",
			10025,
		},
		{
			"smtp://localhost",
			25,
		},
		{
			"smtps://localhost:465",
			465,
		},
		{
			"smtps://localhost:10465",
			10465,
		},
		{
			"smtps://localhost",
			465,
		},
		{
			"smtp+tls://localhost:587",
			587,
		},
		{
			"smtp+tls://localhost:10587",
			10587,
		},
		{
			"smtp+tls://localhost",
			587,
		},
		{
			"enigma://localhost:465",
			465,
		},
		{
			"enigma://localhost:10465",
			10465,
		},
		{
			"enigma://localhost",
			465,
		},
	}

	for k, v := range cases {
		dsn, _ := newDsn(v.dsn)
		if dsn.Port() != v.port {
			t.Errorf("[Case%d] Port: %d (%d)", k, dsn.Port(), v.port)
		}
	}
}

func TestSocket(t *testing.T) {
	cases := []struct {
		dsn    string
		socket string
	}{
		{
			"smtp://localhost:25",
			"localhost:25",
		},
		{
			"smtp://localhost:10025",
			"localhost:10025",
		},
		{
			"smtp://localhost",
			"localhost:25",
		},
		{
			"smtps://localhost:465",
			"localhost:465",
		},
		{
			"smtps://localhost:10465",
			"localhost:10465",
		},
		{
			"smtps://localhost",
			"localhost:465",
		},
		{
			"smtp+tls://localhost:587",
			"localhost:587",
		},
		{
			"smtp+tls://localhost:10587",
			"localhost:10587",
		},
		{
			"smtp+tls://localhost",
			"localhost:587",
		},
		{
			"enigma://localhost:465",
			"localhost:465",
		},
		{
			"enigma://localhost:10465",
			"localhost:10465",
		},
		{
			"enigma://localhost",
			"localhost:465",
		},
	}

	for k, v := range cases {
		dsn, _ := newDsn(v.dsn)
		if dsn.Socket() != v.socket {
			t.Errorf("[Case%d] Socket: %s (%s)", k, dsn.Socket(), v.socket)
		}
	}
}
