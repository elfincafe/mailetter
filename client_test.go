package mailetter

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

var (
	domain   = "example.com"
	mailHost = fmt.Sprintf("mail.%s", domain)
	dsnSmtp  = fmt.Sprintf("smtp://%s", mailHost)
	dsnSmtps = fmt.Sprintf("smtps://%s:465", mailHost)
	dsnTls   = fmt.Sprintf("smtp+tls://%s:25", mailHost)
)

func TestClientNew(t *testing.T) {
	cases := []struct {
		dsn      string
		expected string
	}{
		{
			dsnSmtp,
			"*mailetter.MaiLetter",
		},
		{
			dsnSmtps,
			"*mailetter.MaiLetter",
		},
		{
			dsnTls,
			"*mailetter.MaiLetter",
		},
	}
	for k, v := range cases {
		m := New(v.dsn)
		if reflect.TypeOf(m).String() != v.expected {
			t.Errorf(`[Case%d] %v`, k, reflect.TypeOf(m))
		}
	}
}

func TestMaiLetterLocalName(t *testing.T) {
	cases := []struct {
		call bool
		name string
	}{
		{
			false,
			"localhost",
		},
		{
			true,
			mailHost,
		},
	}

	for k, v := range cases {
		m := New(dsnSmtp)
		if v.call {
			m.LocalName(v.name)
		}
		if m.localName != v.name {
			t.Errorf(`[Case%d] %s (%s)`, k, m.localName, v.name)
		}
	}
}

func TestClientAuthByPlain(t *testing.T) {

}

func TestClientAuthByLogin(t *testing.T) {

}

func TestClientAuthByCramMd5(t *testing.T) {

}

func TestClientHeader(t *testing.T) {

}

func TestClientFrom(t *testing.T) {

}

func TestClientTo(t *testing.T) {

}

func TestClientCc(t *testing.T) {

}

func TestClientBcc(t *testing.T) {

}

func TestClientSubject(t *testing.T) {

}

func TestClientBody(t *testing.T) {

}

func TestClientSet(t *testing.T) {

}

func TestClientSend(t *testing.T) {
	cases := []struct {
		dsn     string
		from    *Address
		to      []*Address
		subject string
		body    io.Reader
		vals    map[string]any
	}{
		{
			dsnSmtps,
			newAddress(fmt.Sprintf("from@%s", domain), "送信者"),
			[]*Address{newAddress(fmt.Sprintf("to+1@%s", domain), "受信者")},
			"テスト件名",
			strings.NewReader("{{.Name}}"),
			map[string]any{"Name": "Mr. Recipient"},
		},
	}
	for k, v := range cases {
		m := newData()
		m.setFrom(v.from)
		for _, t := range v.to {
			m.setTo(t)
		}
		m.setSubject(v.subject)
		m.setBody(v.body)
		for key, val := range v.vals {
			m.setValue(key, val)
		}
		ml := New(v.dsn)
		err := ml.Send()
		if err != nil {
			t.Errorf("[Case%d] %v", k, err)
		}
	}
}

func TestClientReset(t *testing.T) {

}

func TestClientClose(t *testing.T) {
	t.Errorf("Error")
}

func TestClientIsConnect(t *testing.T) {
	cases := []struct {
		dsn       string
		connected bool
	}{
		{
			dsnSmtp,
			false,
		},
		{
			dsnSmtp,
			true,
		},
	}
	k := 0
	v := cases[k]
	ml := New(v.dsn)
	if ml.isConnected() != v.connected {
		t.Errorf("[Case%d] %v(%v)", k, ml.isConnected(), v.connected)
	}
	k = 1
	v = cases[k]
	ml = New(v.dsn)
	dsn := newDsn(v.dsn)
	ml.connectBySmtp(dsn)
	if ml.isConnected() != v.connected {
		t.Errorf("[Case%d] %v(%v)", k, ml.isConnected(), v.connected)
	}
}

func TestClientConnect(t *testing.T) {

}

func TestClientConnectBySmtps(t *testing.T) {
	cases := []struct {
		dsn        string
		clientType string
		errMsg     string
	}{
		{
			"smtps://example.com:",
			"*smtp.Client",
			"",
		},
	}
	for k, v := range cases {
		ml := New(v.dsn)
		dsn := newDsn(v.dsn)
		_, err := ml.connectBySmtps(dsn)
		if err != nil {
			if !strings.Contains(err.Error(), v.errMsg) {
				t.Errorf("[Case%d] %v", k, err)
			}
			continue
		}
		typ := reflect.TypeOf(ml.conn).String()
		if ml.conn != nil && typ != v.clientType {
			t.Errorf("[Case%d] %s (%s)", k, typ, v.clientType)
		}
	}
}

func TestClientConnectBySmtp(t *testing.T) {
	cases := []struct {
		dsn        string
		clientType string
		errMsg     string
	}{
		{
			dsnSmtp,
			"*smtp.Client",
			"",
		},
	}
	for k, v := range cases {
		ml := New(v.dsn)
		dsn := newDsn(v.dsn)
		_, err := ml.connectBySmtp(dsn)
		if err != nil {
			if !strings.Contains(err.Error(), v.errMsg) {
				t.Errorf("[Case%d] %v", k, err)
			}
			continue
		}
		typ := reflect.TypeOf(ml.conn).String()
		if ml.conn != nil && typ != v.clientType {
			t.Errorf("[Case%d] %s (%s)", k, typ, v.clientType)
		}
	}
}

func TestClientConnectWithTls(t *testing.T) {
	cases := []struct {
		dsn        string
		clientType string
		errMsg     string
	}{
		{
			dsnTls,
			"*smtp.Client",
			"",
		},
	}
	for k, v := range cases {
		ml := New(v.dsn)
		dsn := newDsn(v.dsn)
		_, err := ml.connectBySmtp(dsn)
		if err != nil {
			if !strings.Contains(err.Error(), v.errMsg) {
				t.Errorf("[Case%d] %v", k, err)
			}
			continue
		}
		typ := reflect.TypeOf(ml.conn).String()
		if ml.conn != nil && typ != v.clientType {
			t.Errorf("[Case%d] %s (%s)", k, typ, v.clientType)
		}
	}
}
