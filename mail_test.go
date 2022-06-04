package mailetter

import (
	"fmt"
	"testing"
)

func TestTo(t *testing.T) {
	addr1 := "test+1@example.com"
	name1 := "Test1"
	addr2 := "test+2@example.com"
	name2 := "Test2"

	m := NewMail()
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

	m := NewMail()
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

	m := NewMail()
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

	m := NewMail()
	m.Subject(subj1)
	if m.subj != subj1 {
		t.Errorf("Subject:%s(%s)", m.subj, subj1)
	}
	m.Subject(subj2)
	if m.subj != subj2 {
		t.Errorf("Subject:%s(%s)", m.subj, subj2)
	}
}

func TestBody(t *testing.T) {
	body1 := "テスト本文"
	body2 := "TestBody"

	m := NewMail()
	m.Body(body1)
	if m.body != body1 {
		t.Errorf("Body:%s(%s)", m.body, body1)
	}
	m.Body(body2)
	if m.body != body2 {
		t.Errorf("Body:%s(%s)", m.body, body2)
	}
}

func TestCreate(t *testing.T) {
	m := NewMail()
	m.From(NewAddr("from@example.com", "送信者"))
	m.To(NewAddr("to1@example.com", "To受信者1"))
	m.To(NewAddr("to2@example.com", "To受信者2"))
	m.Cc(NewAddr("cc1@example.com", "Cc受信者1"))
	m.Cc(NewAddr("cc2@example.com", "Cc受信者2"))
	m.Bcc(NewAddr("bcc1@example.com", "Bcc受信者1"))
	m.Bcc(NewAddr("bcc2@example.com", "Bcc受信者2"))
	m.Subject("テスト件名 TestSubject")
	m.Body("テスト本文")
	fmt.Println(m.create())
}
