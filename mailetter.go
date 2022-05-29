package mailetter

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/smtp"
	"strings"
	"syscall"
	"time"
)

type MaiLetter struct {
	dsn         *Dsn
	client      *smtp.Client
	hostname    string
	to          []*Addr
	cc          []*Addr
	bcc         []*Addr
	subject     string
	body        string
	from        *Addr
	vars        map[string]interface{}
	attachments []*Attachment
	auth        *Auth
	border      string
}

func New(dsn string, addr *Addr) (*MaiLetter, error) {
	m := new(MaiLetter)
	tmp, err := NewDsn(dsn)
	if err != nil {
		return nil, err
	}
	m.dsn = tmp
	m.client = nil
	m.from = addr
	m.border = border()

	uname := new(syscall.Utsname)
	err = syscall.Uname(uname)
	if err != nil {
		return nil, err
	}
	hostname := strings.Builder{}
	for _, v := range uname.Nodename {
		hostname.WriteString(string(rune(v)))
	}
	m.hostname = hostname.String()

	return m, nil
}

func (m *MaiLetter) Auth(auth *Auth) {
	m.auth = auth
}

func (m *MaiLetter) To(addr *Addr) {
	m.to = append(m.to, addr)
}

func (m *MaiLetter) Cc(addr *Addr) {
	m.cc = append(m.cc, addr)
}

func (m *MaiLetter) Bcc(addr *Addr) {
	m.bcc = append(m.bcc, addr)
}

func (m *MaiLetter) Subject(subject string) {
	m.subject = subject
}

func (m *MaiLetter) Body(body string) {
	m.body = body
}

func (m *MaiLetter) Load(path string) {

}

func (m *MaiLetter) Set(key string, val string) {
	m.vars[key] = val
}

func (m *MaiLetter) Attach(a *Attachment) {
	m.attachments = append(m.attachments, a)
}

func (m *MaiLetter) Send() error {
	var err error
	if !m.isConnected() {
		return m.connect()
	}

	err = m.client.Hello(m.hostname)
	if err != nil {
		return err
	}
	err = m.client.Mail(m.from.Addr())
	for _, addr := range m.to {
		m.client.Rcpt(addr.Addr())
	}
	for _, addr := range m.cc {
		m.client.Rcpt(addr.Addr())
	}
	for _, addr := range m.bcc {
		m.client.Rcpt(addr.Addr())
	}

	return nil
}

func (m *MaiLetter) Reset() error {
	err := m.client.Reset()
	return err
}

func (m *MaiLetter) Noop() error {
	err := m.client.Noop()
	return err
}

func (m *MaiLetter) Quit() error {
	err := m.client.Quit()
	return err
}

func (m *MaiLetter) isConnected() bool {
	if m.client != nil {
		return true
	} else {
		return false
	}
}

func (m *MaiLetter) connect() error {
	fmt.Println(m.dsn.String())
	client, err := smtp.Dial(m.dsn.String())
	fmt.Println(err)
	if err != nil {
		return err
	}
	m.client = client

	return nil
}

func encode(s string) string {
	needsEnc := false
	for _, v := range s {
		if v > 127 {
			needsEnc = true
			break
		}
	}
	if !needsEnc {
		return s
	}
	return fmt.Sprintf("=?UTF-8?%s?=", base64.StdEncoding.EncodeToString([]byte(s)))
}

func encodeBinary(content []byte) string {
	return base64.StdEncoding.EncodeToString(content)
}

func border() string {
	length := 24
	s := []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
		"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
		"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	}
	l := len(s)
	rand.Seed(time.Now().UnixNano())
	sb := strings.Builder{}
	sb.WriteString(strings.Repeat("-", 12))
	for i := 0; i < length; i++ {
		idx := rand.Intn(l - 1)
		sb.WriteString(s[idx])
	}
	return sb.String()
}

func getResponse (r *bufio.Reader) {

}
