package mailetter

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
	"text/template"
)

func TestNewMail(t *testing.T) {
	from := NewAddr("from@example.com", "Sender")
	m := NewMail(from)
	typ := reflect.TypeOf(m)
	expect := reflect.TypeOf((*Mail)(nil)).String()
	if typ.String() != expect {
		t.Errorf("[Case%d] Type: %s != %s", 0, typ.String(), expect)
	}
	if m.from != from {
		t.Errorf("[Case%d] From: %v != %v", 1, m.from, from)
	}
}

func TestMailReset(t *testing.T) {
	cases := []struct {
		headers map[string]header
		to      []*Addr
		cc      []*Addr
		bcc     []*Addr
		subject *template.Template
		body    *template.Template
		vars    map[string]any
	}{
		{
			map[string]header{"test-header": header{"Test-Header", "Test Header Value"}},
			[]*Addr{NewAddr("To Address", "to@example.com")},
			[]*Addr{NewAddr("Cc Address", "cc@example.com")},
			[]*Addr{NewAddr("Bcc Address", "bcc@example.com")},
			template.Must(template.New("Subject").Parse("Dear. {{.Name}}")),
			template.Must(template.New("Body").Parse("")),
			map[string]any{"Name": "Example User"},
		},
	}
	from := NewAddr("from@example.com", "From Address")
	m := NewMail(from)
	for k, v := range cases {
		m.headers = v.headers
		m.to = v.to
		m.cc = v.cc
		m.bcc = v.bcc
		m.subject = v.subject
		m.body = v.body
		m.vars = v.vars
		m.Reset()
		if len(m.headers) != 0 || len(m.to) != 0 || len(m.cc) != 0 || len(m.bcc) != 0 || m.subject != nil || m.body != nil || len(m.vars) != 0 {
			t.Errorf("[Case%d] Header:%v, To:%v, Cc:%v, Bcc:%v, Subject:%v, Body:%v, Vars:%v", k, m.headers, m.to, m.cc, m.bcc, m.subject, m.body, m.vars)
		}

	}
}

func TestMailHeader(t *testing.T) {
	cases := []struct {
		key    string
		val    string
		exists bool
	}{
		{"Date", "Tue, 5 Mar 2024 21:53:04 +0900", false},
		{"From", "Test From <from@example.com>", false},
		{"To", "Test To <to@example.com>", false},
		{"Cc", "Test Cc <cc@example.com>", false},
		{"Bcc", "bcc@example.com", false},
		{"Subject", "Test Subject", false},
		{"Reply-To", "reply-to@example.com", false},
		{"Return-Path", "return-path@example.com", false},
		{"X-Mailer", "Test MTU 1", true},
		{"Message-ID", "<1234567890ABCDEFGHIJKLMN@example.com>", true},
		{"X-Mailer", "Test MTU 2", true},
	}
	count := 2
	from := NewAddr("from@example.com", "Sender")
	m := NewMail(from)
	for k, v := range cases {
		m.Header(v.key, v.val)
		key := strings.ToLower(v.key)
		_, ok := m.headers[key]
		if ok == v.exists {
			if !v.exists {
				continue
			}
		} else {
			t.Errorf(`[Case%d] "%s" doesn't exist.`, k, v.key)
		}
		if m.headers[key].key != v.key {
			t.Errorf(`[Case%d] Key: %s (%s)`, k, m.headers[key].key, v.key)
		}
		if m.headers[key].val != v.val {
			t.Errorf(`[Case%d] Value: %s (%s)`, k, m.headers[key].val, v.val)
		}
	}
	if len(m.headers) != count {
		t.Errorf("[Case%d] Count: %d (%d)", 999, len(m.headers), count)
	}
}

func TestMailTo(t *testing.T) {
	cases := []struct {
		addr string
	}{
		{"to+0@example.com"},
		{"to+1@example.com"},
		{"to+2@example.com"},
	}
	from := NewAddr("from@example.com", "Sender")
	m := NewMail(from)
	for k, v := range cases {
		a := NewAddr(v.addr, "")
		m.To(a)
		if m.to[k] != a {
			t.Errorf("[Case%d] Address: %v (%v)", k, m.to[k], a)
		}
	}
	if len(m.to) != len(cases) {
		t.Errorf("[Case%d] Count: %d (%d)", 999, len(m.to), len(cases))
	}
}

func TestMailCc(t *testing.T) {
	cases := []struct {
		addr string
	}{
		{"cc+0@example.com"},
		{"cc+1@example.com"},
		{"cc+2@example.com"},
	}
	from := NewAddr("from@example.com", "Sender")
	m := NewMail(from)
	for k, v := range cases {
		a := NewAddr(v.addr, "")
		m.Cc(a)
		if m.cc[k] != a {
			t.Errorf("[Case%d] Address: %v (%v)", k, m.cc[k], a)
		}
	}
	if len(m.cc) != len(cases) {
		t.Errorf("[Case%d] Count: %d (%d)", 999, len(m.cc), len(cases))
	}
}

func TestMailBcc(t *testing.T) {
	cases := []struct {
		addr string
	}{
		{
			"bcc+0@example.com",
		},
		{
			"bcc+1@example.com",
		},
		{
			"bcc+2@example.com",
		},
	}
	from := NewAddr("from@example.com", "Sender")
	m := NewMail(from)
	for k, v := range cases {
		a := NewAddr(v.addr, "")
		m.Bcc(a)
		if m.bcc[k] != a {
			t.Errorf("[Case%d] Address: %v (%v)", k, m.bcc[k], a)
		}
	}
	if len(m.bcc) != len(cases) {
		t.Errorf("[Case%d] Count: %d (%d)", 999, len(m.bcc), len(cases))
	}
}

func TestMailReplyTo(t *testing.T) {
	cases := []struct {
		from  *Addr
		reply *Addr
	}{
		{
			NewAddr("from@example.com", ""),
			NewAddr("reply-to@example.com", ""),
		},
	}

	for k, v := range cases {
		m := NewMail(v.from)
		if m.replyTo.addr != v.from.addr {
			t.Errorf(`[Case%d] InitAddr: %s (%s)`, k, m.replyTo.addr, v.from.addr)
			continue
		}
		m.ReplyTo(v.reply)
		if m.replyTo.addr != v.reply.addr {
			t.Errorf(`[Case%d] InitAddr: %s (%s)`, k, m.replyTo.addr, v.reply.addr)
		}
	}
}

func TestMailReturnPath(t *testing.T) {
	cases := []struct {
		from *Addr
		ret  *Addr
	}{
		{
			NewAddr("from@example.com", ""),
			NewAddr("return-path@example.com", ""),
		},
	}

	for k, v := range cases {
		m := NewMail(v.from)
		if m.replyTo.addr != v.from.addr {
			t.Errorf(`[Case%d] InitAddr: %s (%s)`, k, m.replyTo.addr, v.from.addr)
			continue
		}
		m.ReplyTo(v.ret)
		if m.replyTo.addr != v.ret.addr {
			t.Errorf(`[Case%d] InitAddr: %s (%s)`, k, m.replyTo.addr, v.ret.addr)
		}
	}
}

func TestMailSubject(t *testing.T) {
	cases := []struct {
		subject string
		vars    any
	}{
		{
			"Subject1",
			nil,
		},
		{
			"Dear {{.Name}}",
			nil,
		},
	}

	from := NewAddr("from@example.com", "Sender")
	m := NewMail(from)
	for k, v := range cases {
		m.Subject(v.subject)
		typ := reflect.TypeOf(m.subject).String()
		if typ != "*template.Template" {
			t.Errorf(`[Case%d] %v`, k, m)
		}
	}
}

func TestMailBody(t *testing.T) {
	cases := []struct {
		body string
		vars any
	}{
		{
			"Test Body Part1",
			nil,
		},
		{
			"Test\r\nBody\r\nPart2\r\n{{.Name}}",
			nil,
		},
	}

	from := NewAddr("from@example.com", "Sender")
	m := NewMail(from)
	for k, v := range cases {
		m.Subject(v.body)
		typ := reflect.TypeOf(m.subject).String()
		fmt.Println(typ)
		if typ != "*template.Template" {
			t.Errorf(`[Case%d] %v`, k, m)
		}
	}
}

func TestMailSet(t *testing.T) {
	cases := []struct {
		key string
		val any
	}{
		{
			"a",
			"abc",
		},
		{
			"b",
			1,
		},
		{
			"a",
			[]string{"a", "b", "c"},
		},
	}
	expected := 2

	from := NewAddr("from@example.com", "Sender")
	m := NewMail(from)
	for k, v := range cases {
		m.Set(v.key, v.val)
		typ1 := reflect.TypeOf(m.vars[v.key])
		typ2 := reflect.TypeOf(v.val)
		if typ1 != typ2 {
			t.Errorf("[Case%d] Type: %v (%v)", k, typ1, typ2)
		}
	}
	if len(m.vars) != expected {
		t.Errorf("[Case%d] Count: %d (%d)", 999, len(m.vars), expected)
	}
}

func TestMailString(t *testing.T) {
	cases := []struct {
		headers map[string]header
		to      []*Addr
		cc      []*Addr
		bcc     []*Addr
		subject string
		body    io.Reader
		vars    map[string]any
	}{
		{
			map[string]header{"test-header": {"Test-Header", "Test Header Value"}},
			[]*Addr{NewAddr("to0@example.com", "受信者To0"), NewAddr("to1@example.com", "受信者To1")},
			[]*Addr{NewAddr("cc0@example.com", "受信者Cc0"), NewAddr("cc1@example.com", "受信者Cc1"), NewAddr("cc2@example.com", "受信者Cc2")},
			[]*Addr{NewAddr("bcc0@example.com", "受信者Bcc0")},
			"[{{.ServiceName}}] {{.Name}}様 新商品のお知らせ",
			strings.NewReader("{{.Name}}様\nいつもご利用ありがとうございます。\n{{.ServiceName}}カスタマーサポートでございます。"),
			map[string]any{"ServiceName": "ECサービス", "Name": "ECサービスユーザー"},
		},
	}

	m := NewMail(NewAddr("from@example.com", "送信者"))
	for i, c := range cases {
		for _, v := range c.headers {
			m.Header(v.key, v.val)
		}
		for _, v := range c.to {
			m.To(v)
		}
		for _, v := range c.cc {
			m.Cc(v)
		}
		for _, v := range c.bcc {
			m.Bcc(v)
		}
		for k, v := range c.vars {
			m.Set(k, v)
		}
		m.Subject(c.subject)
		m.Body(c.body)
		fmt.Println(m.String())

		if false {
			t.Errorf(`[Case%d]`, i)
		}
	}
}
