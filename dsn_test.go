package mailetter

import (
	"testing"
)

func TestNewDsn(t *testing.T) {
	var dsn *Dsn
	// var err error
	var str string

	str = "smtp://example.com"
	dsn, _ = NewDsn(str)
	if dsn.scheme != "smtp" || dsn.host != "example.com" || dsn.port != 587 {
		t.Errorf("DSN: %s, Scheme:%s, Host:%s, Port:%d)", str, dsn.scheme, dsn.host, dsn.port)
	}
	str = "smtps://example.com:465"
	dsn, _ = NewDsn(str)
	if dsn.scheme != "smtps" || dsn.host != "example.com" || dsn.port != 465 {
		t.Errorf("DSN: %s, Scheme:%s, Host:%s, Port:%d)", str, dsn.scheme, dsn.host, dsn.port)
	}
	str = "smtp+tls://example.com:2525"
	dsn, _ = NewDsn(str)
	if dsn.scheme != "smtp+tls" || dsn.host != "example.com" || dsn.port != 2525 {
		t.Errorf("DSN: %s, Scheme:%s, Host:%s, Port:%d)", str, dsn.scheme, dsn.host, dsn.port)
	}
}
