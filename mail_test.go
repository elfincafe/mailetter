package mailetter

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestNewMail(t *testing.T) {
	from, _ := NewAddress("from@example.com", "Sender")
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
	a, _ := NewAddress()
	cases := []struct {
		headers map[string]header
		to      []*Address
		cc      []*Address
		bcc     []*Address
		subject *template.Template
		body    *template.Templte
		vars    map[string]any
	}{
		{},
	}
	for k, v := range cases {
	}
}

func TestMailHeader(t *testing.T) {
	cases := []struct {
		key string
		val string
	}{
		{"Subject", "Mail Subject"},
		{"X-Mailer", "Test MTU 1"},
		{"Message-ID", "<1234567890ABCDEFGHIJKLMN@example.com>"},
		{"X-Mailer", "Test MTU 2"},
	}
	expected := 3
	from, _ := NewAddress("from@example.com", "Sender")
	m := NewMail(from)
	for k, v := range cases {
		m.Header(v.key, v.val)
		key := strings.ToLower(v.key)
		if _, ok := m.headers[key]; !ok {
			t.Errorf(`[Case%d] "%s" doesn't exist.`, k, v.key)
		}
		if m.headers[key].key != v.key {
			t.Errorf(`[Case%d] Key: %s (%s)`, k, m.headers[key].key, v.key)
		}
		if m.headers[key].val != v.val {
			t.Errorf(`[Case%d] Value: %s (%s)`, k, m.headers[key].val, v.val)
		}
	}
	if len(m.headers) != expected {
		t.Errorf("[Case%d] Count: %d (%d)", 999, len(m.headers), expected)
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
	from, _ := NewAddress("from@example.com", "Sender")
	m := NewMail(from)
	for k, v := range cases {
		a, _ := NewAddress(v.addr, "")
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
	from, _ := NewAddress("from@example.com", "Sender")
	m := NewMail(from)
	for k, v := range cases {
		a, _ := NewAddress(v.addr, "")
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
		{"bcc+0@example.com"},
		{"bcc+1@example.com"},
		{"bcc+2@example.com"},
	}
	from, _ := NewAddress("from@example.com", "Sender")
	m := NewMail(from)
	for k, v := range cases {
		a, _ := NewAddress(v.addr, "")
		m.Bcc(a)
		if m.bcc[k] != a {
			t.Errorf("[Case%d] Address: %v (%v)", k, m.bcc[k], a)
		}
	}
	if len(m.bcc) != len(cases) {
		t.Errorf("[Case%d] Count: %d (%d)", 999, len(m.bcc), len(cases))
	}
}

func TestMailSubject(t *testing.T) {
	cases := []struct {
		subject string
		vars    any
	}{
		{"Subject1", nil},
		{"Dear {{.Name}}", nil},
	}

	from, _ := NewAddress("from@example.com", "Sender")
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
		{"Test Body Part1", nil},
		{"Test\r\nBody\r\nPart2\r\n{{.Name}}", nil},
	}

	from, _ := NewAddress("from@example.com", "Sender")
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
		{"a", "abc"},
		{"b", 1},
		{"a", []string{"a", "b", "c"}},
	}
	expected := 2

	from, _ := NewAddress("from@example.com", "Sender")
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

	from, _ := NewAddress("from@example.com", "Sender")
	m := NewMail(from)
	for _, v := range []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9} {
		to, _ := NewAddress(fmt.Sprintf("to+%d@example.com", v), "受信者1")
		m.To(to)
	}
	m.Subject("貨表示を円貨にした場合、平均取得価額は国内約定日の10時30分までは参考レートに為替掛目1％を加えて計算している")
	m.Body("通貨表示を円貨とした場合の時価評価額・評価損益は、現在の参考為替レートを利用して円換算額を算出しております。\r\n従って、実際の円貨決済による売却時の円換算レート(TTB)、\r\nまたそれにより算出される売却価額 (売却時の費用を考慮しない) とは異なります。")
	t.Errorf("%s", m.String())
}
