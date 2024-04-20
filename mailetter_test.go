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

func TestMaiLetterNew(t *testing.T) {
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
		m, _ := New(v.dsn)
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
		m, _ := New(dsnSmtp)
		if v.call {
			m.LocalName(v.name)
		}
		if m.localName != v.name {
			t.Errorf(`[Case%d] %s (%s)`, k, m.localName, v.name)
		}
	}
}

func TestMaiLetterConnectWithTls(t *testing.T) {
	cases := []struct {
		dsn        string
		clientType string
		errMsg     string
	}{
		{
			dsnSmtps,
			"*smtp.Client",
			"",
		},
	}
	for k, v := range cases {
		ml, _ := New(v.dsn)
		err := ml.connectWithSsl()
		if err != nil {
			if !strings.Contains(err.Error(), v.errMsg) {
				t.Errorf("[Case%d] %v", k, err)
			}
			continue
		}
		typ := reflect.TypeOf(ml.client).String()
		if ml.client != nil && typ != v.clientType {
			t.Errorf("[Case%d] %s (%s)", k, typ, v.clientType)
		}
	}
}

func TestMaiLetterConnectWithoutTls(t *testing.T) {
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
		ml, _ := New(v.dsn)
		err := ml.connectWithoutSsl()
		if err != nil {
			if !strings.Contains(err.Error(), v.errMsg) {
				t.Errorf("[Case%d] %v", k, err)
			}
			continue
		}
		typ := reflect.TypeOf(ml.client).String()
		if ml.client != nil && typ != v.clientType {
			t.Errorf("[Case%d] %s (%s)", k, typ, v.clientType)
		}
	}
}

func TestMaiLetterConnectAndStartTls(t *testing.T) {
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
		ml, _ := New(v.dsn)
		err := ml.connectWithoutSsl()
		if err != nil {
			if !strings.Contains(err.Error(), v.errMsg) {
				t.Errorf("[Case%d] %v", k, err)
			}
			continue
		}
		typ := reflect.TypeOf(ml.client).String()
		if ml.client != nil && typ != v.clientType {
			t.Errorf("[Case%d] %s (%s)", k, typ, v.clientType)
		}
	}
}

func TestMaiLetterIsConnected(t *testing.T) {
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
	ml, _ := New(v.dsn)
	if ml.isConnected() != v.connected {
		t.Errorf("[Case%d] %v(%v)", k, ml.isConnected(), v.connected)
	}
	k = 1
	v = cases[k]
	ml, _ = New(v.dsn)
	ml.connectWithoutSsl()
	if ml.isConnected() != v.connected {
		t.Errorf("[Case%d] %v(%v)", k, ml.isConnected(), v.connected)
	}
}

func TestMaiLetterSend(t *testing.T) {
	cases := []struct {
		dsn     string
		from    *Addr
		to      []*Addr
		subject string
		body    io.Reader
		vals    map[string]any
	}{
		{
			dsnSmtps,
			NewAddr(fmt.Sprintf("from@%s", domain), "送信者"),
			[]*Addr{NewAddr(fmt.Sprintf("to+1@%s", domain), "受信者")},
			"テスト件名",
			strings.NewReader("{{.Name}}"),
			map[string]any{"Name": "Mr. Recipient"},
		},
	}
	for k, v := range cases {
		m := NewMail(v.from)
		for _, t := range v.to {
			m.To(t)
		}
		m.Subject(v.subject)
		m.Body(v.body)
		for key, val := range v.vals {
			m.Set(key, val)
		}
		ml, _ := New(v.dsn)
		err := ml.Send(m)
		if err != nil {
			t.Errorf("[Case%d] %v", k, err)
		}
	}
}
