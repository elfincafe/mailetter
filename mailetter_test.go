package mailetter

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

const (
	DSN       = "smtp://smtp.example.com:25"
	FROM_ADDR = "test@example.com"
	FROM_NAME = ""
)

func TestIsConnected(t *testing.T) {
	goodDsn := DSN
	m, _ := New(goodDsn, NewAddr(FROM_ADDR, FROM_NAME))
	m.connect()
	fmt.Println(m.client)
	if !m.isConnected() {
		t.Errorf("Connection Failed. Something is wrong.")
	}
	badDsn := strings.ReplaceAll(DSN, ":25", ":12345")
	m, _ = New(badDsn, NewAddr(FROM_ADDR, FROM_NAME))
	m.connect()
	if m.isConnected() {
		t.Errorf("Connection Success. Something is wrong.")
	}
}

func TestConnect(t *testing.T) {
	_, err := New(DSN, NewAddr(FROM_ADDR, FROM_NAME))
	if err != nil {
		fmt.Println(err)
	}
}

func TestTo(t *testing.T) {
	addr1 := "test+1@example.com"
	name1 := "Test1"
	addr2 := "test+2@example.com"
	name2 := "Test2"

	m, _ := New(DSN, NewAddr(FROM_ADDR, FROM_NAME))
	m.To(NewAddr(addr1, name1))
	if len(m.to) != 1 || m.to[0].Addr() != addr1 || m.to[0].Name() != name1 {
		t.Errorf("Count=%d(%d), Address=%s(%s), Name=%s(%s)", len(m.to), 1, m.to[0].Addr(), addr1, m.to[0].Name(), name1)
	}
	m.To(NewAddr(addr2, name2))
	if len(m.to) != 2 || m.to[1].Addr() != addr2 || m.to[1].Name() != name2 {
		t.Errorf("Count=%d(%d), Address=%s(%s), Name=%s(%s)", len(m.to), 2, m.to[1].Addr(), addr2, m.to[1].Name(), name2)
	}
}

func TestCc(t *testing.T) {
	addr1 := "test+1@example.com"
	name1 := "Test1"
	addr2 := "test+2@example.com"
	name2 := "Test2"

	m, _ := New(DSN, NewAddr(FROM_ADDR, FROM_NAME))
	m.Cc(NewAddr(addr1, name1))
	if len(m.cc) != 1 || m.cc[0].Addr() != addr1 || m.cc[0].Name() != name1 {
		t.Errorf("Count=%d(%d), Address=%s(%s), Name=%s(%s)", len(m.cc), 1, m.cc[0].Addr(), addr1, m.cc[0].Name(), name1)
	}
	m.Cc(NewAddr(addr2, name2))
	if len(m.cc) != 2 || m.cc[1].Addr() != addr2 || m.cc[1].Name() != name2 {
		t.Errorf("Count=%d(%d), Address=%s(%s), Name=%s(%s)", len(m.cc), 2, m.cc[1].Addr(), addr2, m.cc[1].Name(), name2)
	}
}

func TestBcc(t *testing.T) {
	addr1 := "test+1@example.com"
	addr2 := "test+2@example.com"

	m, _ := New(DSN, NewAddr(FROM_ADDR, FROM_NAME))
	m.Bcc(NewAddr(addr1, ""))
	if len(m.bcc) != 1 || m.bcc[0].Addr() != addr1 {
		t.Errorf("Count=%d(%d), Address=%s(%s)", len(m.bcc), 1, m.bcc[0].Addr(), addr1)
	}
	m.Bcc(NewAddr(addr2, ""))
	if len(m.bcc) != 2 || m.bcc[1].Addr() != addr2 {
		t.Errorf("Count=%d(%d), Address=%s(%s)", len(m.bcc), 2, m.bcc[1].Addr(), addr2)
	}
}

func TestSubject(t *testing.T) {
	subj1 := "テスト件名"
	subj2 := "TestSubject"
	m, _ := New(DSN, NewAddr(FROM_ADDR, FROM_NAME))
	m.Subject(subj1)
	if m.subject != subj1 {
		t.Errorf("Subject:%s(%s)", m.subject, subj1)
	}
	m.Subject(subj2)
	if m.subject != subj2 {
		t.Errorf("Subject:%s(%s)", m.subject, subj2)
	}
}

func TestBody(t *testing.T) {
	body1 := "テスト本文"
	body2 := "TestBody"
	m, _ := New(DSN, NewAddr(FROM_ADDR, FROM_NAME))
	m.Body(body1)
	if m.body != body1 {
		t.Errorf("Body:%s(%s)", m.body, body1)
	}
	m.Body(body2)
	if m.body != body2 {
		t.Errorf("Body:%s(%s)", m.body, body2)
	}
}

func TestBorder(t *testing.T) {
	border := border()
	re := regexp.MustCompile(`-{12}[0-9a-zA-Z]{24}`)
	if !re.MatchString(border) {
		t.Errorf(`Invalid Border %s`, border)
	}
}
