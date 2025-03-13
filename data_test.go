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
	if data.headers == nil || len(data.headers) != 0 || data.hdrOrder == nil || len(data.hdrOrder) != 0 {
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
	data.reset()
	if data.headers == nil || len(data.headers) != 0 || data.hdrOrder == nil || len(data.hdrOrder) != 0 {
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

func TestDataSetHeader(t *testing.T) {
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
	m := newData()
	from := newAddress("from@example.com", "送信者")
	m.setFrom(from)
	for k, v := range cases {
		m.setHeader(v.key, v.val)
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
		if m.headers[key].value != v.val {
			t.Errorf(`[Case%d] Value: %s (%s)`, k, m.headers[key].value, v.val)
		}
	}
	if len(m.headers) != count {
		t.Errorf("[Case%d] Count: %d (%d)", 999, len(m.headers), count)
	}
}

func TestDataSetFrom(t *testing.T) {

}

func TestDataSetTo(t *testing.T) {
	cases := []struct {
		addr string
	}{
		{"to+0@example.com"},
		{"to+1@example.com"},
		{"to+2@example.com"},
	}
	m := newData()
	from := newAddress("from@example.com", "送信者")
	m.setFrom(from)
	for k, v := range cases {
		a := newAddress(v.addr, "")
		m.setTo(a)
		if m.to[k] != a {
			t.Errorf("[Case%d] Address: %v (%v)", k, m.to[k], a)
		}
	}
	if len(m.to) != len(cases) {
		t.Errorf("[Case%d] Count: %d (%d)", 999, len(m.to), len(cases))
	}
}

func TestDataSetCc(t *testing.T) {
	cases := []struct {
		addr string
	}{
		{"cc+0@example.com"},
		{"cc+1@example.com"},
		{"cc+2@example.com"},
	}
	m := newData()
	from := newAddress("from@example.com", "Sender")
	m.setFrom(from)
	for k, v := range cases {
		a := newAddress(v.addr, "")
		m.setCc(a)
		if m.cc[k] != a {
			t.Errorf("[Case%d] Address: %v (%v)", k, m.cc[k], a)
		}
	}
	if len(m.cc) != len(cases) {
		t.Errorf("[Case%d] Count: %d (%d)", 999, len(m.cc), len(cases))
	}
}

func TestDataSetBcc(t *testing.T) {
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
	m := newData()
	from := newAddress("from@example.com", "Sender")
	m.setFrom(from)
	for k, v := range cases {
		a := newAddress(v.addr, "")
		m.setBcc(a)
		if m.bcc[k] != a {
			t.Errorf("[Case%d] Address: %v (%v)", k, m.bcc[k], a)
		}
	}
	if len(m.bcc) != len(cases) {
		t.Errorf("[Case%d] Count: %d (%d)", 999, len(m.bcc), len(cases))
	}
}

func TestDataSetReturnPath(t *testing.T) {
	cases := []struct {
		from *Address
		ret  *Address
	}{
		{
			newAddress("from@example.com", ""),
			newAddress("return-path@example.com", ""),
		},
	}

	for k, v := range cases {
		m := newData()
		if m.replyTo.address != v.from.address {
			t.Errorf(`[Case%d] InitAddr: %s (%s)`, k, m.replyTo.address, v.from.address)
			continue
		}
		m.setReplyTo(v.ret)
		if m.replyTo.address != v.ret.address {
			t.Errorf(`[Case%d] InitAddr: %s (%s)`, k, m.replyTo.address, v.ret.address)
		}
	}
}

func TestDataSetReplyTo(t *testing.T) {
	cases := []struct {
		from  *Address
		reply *Address
	}{
		{
			newAddress("from@example.com", ""),
			newAddress("reply-to@example.com", ""),
		},
	}

	for k, v := range cases {
		m := newData()
		if m.replyTo.address != v.from.address {
			t.Errorf(`[Case%d] InitAddr: %s (%s)`, k, m.replyTo.address, v.from.address)
			continue
		}
		m.setReplyTo(v.reply)
		if m.replyTo.address != v.reply.address {
			t.Errorf(`[Case%d] InitAddr: %s (%s)`, k, m.replyTo.address, v.reply.address)
		}
	}
}

func TestDataSetSubject(t *testing.T) {
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

	m := newData()
	from := newAddress("from@example.com", "送信者")
	m.setFrom(from)
	for k, v := range cases {
		m.setSubject(v.subject)
		typ := reflect.TypeOf(m.subject).String()
		if typ != "*template.Template" {
			t.Errorf(`[Case%d] %v`, k, m)
		}
	}
}

func TestDataSetBody(t *testing.T) {
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

	m := newData()
	from := newAddress("from@example.com", "送信者")
	m.setFrom(from)
	for k, v := range cases {
		m.setSubject(v.body)
		typ := reflect.TypeOf(m.subject).String()
		fmt.Println(typ)
		if typ != "*template.Template" {
			t.Errorf(`[Case%d] %v`, k, m)
		}
	}
}

func TestDataSetValue(t *testing.T) {
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

	m := newData()
	from := newAddress("from@example.com", "送信者")
	m.setFrom(from)
	for k, v := range cases {
		m.setValue(v.key, v.val)
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

func TestDataString(t *testing.T) {
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
			[]*Address{newAddress("to0@example.com", "受信者To0"), newAddress("to1@example.com", "受信者To1")},
			[]*Address{newAddress("cc0@example.com", "受信者Cc0"), newAddress("cc1@example.com", "受信者Cc1"), newAddress("cc2@example.com", "受信者Cc2")},
			[]*Address{newAddress("bcc0@example.com", "受信者Bcc0")},
			"[{{.ServiceName}}] {{.Name}}様 新商品のお知らせ",
			strings.NewReader("{{.Name}}様\nいつもご利用ありがとうございます。\n{{.ServiceName}}カスタマーサポートでございます。"),
			map[string]any{"ServiceName": "ECサービス", "Name": "ECサービスユーザー"},
		},
	}

	m := newData()
	from := newAddress("from@example.com", "送信者")
	m.setFrom(from)
	for i, c := range cases {
		for _, v := range c.headers {
			m.setHeader(v.key, v.value)
		}
		for _, v := range c.to {
			m.setTo(v)
		}
		for _, v := range c.cc {
			m.setCc(v)
		}
		for _, v := range c.bcc {
			m.setBcc(v)
		}
		for k, v := range c.vars {
			m.setValue(k, v)
		}
		m.setSubject(c.subject)
		m.setBody(c.body)
		_, _ = m.Create()

		if false {
			t.Errorf(`[Case%d]`, i)
		}
	}
}
