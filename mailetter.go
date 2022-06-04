package mailetter

import (
	"fmt"
	"net/smtp"
	"strings"
	"syscall"
)

type MaiLetter struct {
	dsn      *Dsn
	client   *smtp.Client
	hostname string
	vars     map[string]interface{}
	auth     *Auth
	mail     *Mail
}

func New(dsn string, opts map[string]interface{}) (*MaiLetter, error) {
	m := new(MaiLetter)
	tmp, err := NewDsn(dsn)
	if err != nil {
		return nil, err
	}
	m.dsn = tmp
	m.client = nil
	m.auth = nil

	uname := new(syscall.Utsname)
	err = syscall.Uname(uname)
	if err != nil {
		return nil, err
	}
	hostname := strings.Builder{}
	if _, ok := opts["hostname"]; ok {
		hostname.WriteString(opts["hostname"].(string))
	}
	for _, v := range uname.Nodename {
		hostname.WriteString(string(rune(v)))
	}
	m.hostname = hostname.String()

	m.mail = NewMail()

	return m, nil
}

func (m *MaiLetter) Auth(auth *Auth) {
	m.auth = auth
}

func (m *MaiLetter) Mail() *Mail {
	return m.mail
}

func (m *MaiLetter) Send() error {
	var err error
	if !m.isConnected() {
		m.connect()
	}

	err = m.client.Hello(m.hostname)
	if err != nil {
		return err
	}
	m.mail.create()
	return nil
}

func (m *MaiLetter) Reset() error {
	err := m.client.Reset()
	if err != nil {
		return err
	}
	m.mail = NewMail()
	return nil
}

func (m *MaiLetter) Noop() error {
	return m.client.Noop()
}

func (m *MaiLetter) Quit() error {
	return m.client.Quit()
}

func (m *MaiLetter) Close() error {
	return m.client.Close()
}

func (m *MaiLetter) isConnected() bool {
	if m.client != nil {
		return true
	} else {
		return false
	}
}

func (m *MaiLetter) connect() error {
	// fmt.Println(m.dsn.String())
	client, err := smtp.Dial(m.dsn.String())
	if err != nil {
		fmt.Println(err)
		return err
	}
	m.client = client

	return nil
}
