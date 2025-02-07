package mailetter

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
	"text/template"
)

func TestNewData(t *testing.T) {
	data := newData()
	if data == nil {
		t.Errorf(`Data struct is empty`)
	}
	if len(data.headers) == 0 || len(data.hdrOrder) == 0 {
		t.Errorf(`Data struct headers should be empty (Headers: %d, Orders: %d)`, len(data.headers), len(data.hdrOrder))
	}
	if data.from != nil {
		t.Errorf(`"From" should be empty %v`, data.from)
	}
	if data.returnPath != nil {
		t.Errorf(`"Return-Path" should be empty %v`, data.returnPath)
	}
	if data.replyTo != nil {
		t.Errorf(`"Reply-To" should be empty %v`, data.replyTo)
	}
	if len(data.to) != 0 {
		t.Errorf(`"To" should be empty %v`, data.to)
	}
	if len(data.cc) != 0 {
		t.Errorf(`"Cc" should be empty %v`, data.cc)
	}
	if len(data.bcc) != 0 {
		t.Errorf(`"Bcc" should be empty %v`, data.bcc)
	}
	if data.subject != nil {
		t.Errorf(`"Subject" should be empty %v`, data.subject)
	}
	if data.body != nil {
		t.Errorf(`Body should be empty %v`, data.body)
	}
	if len(data.vars) != 0 {
		t.Errorf(`Variables should be empty %v`, data.vars)
	}
}

func TestDataReset(t *testing.T) {
	data := newData()
	data.headers["xtest"] = header{key: "X-Test", value: "X Test Value"}
	data.hdrOrder = append(data.hdrOrder, "xtest")
	data.from = newAddress("from@example.com", "FromAddress")
	data.to = append(data.to, newAddress("to@example.com", "ToAddress"))
	data.cc = append(data.cc, newAddress("cc@example.com", "CcAddress"))
	data.bcc = append(data.bcc, newAddress("bcc@example.com", "BccAddress"))
	data.subject = template.Must(template.New("Subject").Parse("Test Subject"))
	data.body = template.Must(template.New("Body").Parse("Test Body"))
	data.vars["Test"] = "test value"
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
		from  *Address
		reply *Address
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
		from *Address
		ret  *Address
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
		to      []*Address
		cc      []*Address
		bcc     []*Address
		subject string
		body    io.Reader
		vars    map[string]any
	}{
		{
			map[string]header{"test-header": {"Test-Header", "Test Header Value"}},
			[]*Address{NewAddr("to0@example.com", "受信者To0"), NewAddr("to1@example.com", "受信者To1")},
			[]*Address{NewAddr("cc0@example.com", "受信者Cc0"), NewAddr("cc1@example.com", "受信者Cc1"), NewAddr("cc2@example.com", "受信者Cc2")},
			[]*Address{NewAddr("bcc0@example.com", "受信者Bcc0")},
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
