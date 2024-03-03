package mailetter

import (
	"net/smtp"
	"strings"
	"syscall"
)

const (
	br          = "\r\n"
	white_space = " \r\n\t\v\b"
	should_br   = 78
	must_br     = 998
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
	ml := new(MaiLetter)
	tmp, err := NewDsn(dsn)
	if err != nil {
		return nil, err
	}
	ml.dsn = tmp
	ml.client = nil
	ml.auth = nil

	hostname := strings.Builder{}
	if _, ok := opts["hostname"]; ok {
		hostname.WriteString(opts["hostname"].(string))
	} else {
		uname := new(syscall.Utsname)
		err = syscall.Uname(uname)
		if err != nil {
			return nil, err
		}
		for _, v := range uname.Nodename {
			hostname.WriteString(string(rune(v)))
		}
	}
	ml.hostname = hostname.String()

	// ml.mail = NewMail()

	return ml, nil
}

func (ml *MaiLetter) Auth(auth *Auth) {
	ml.auth = auth
}

func (ml *MaiLetter) Mail() *Mail {
	return ml.mail
}

func (ml *MaiLetter) Send() error {

	var err error
	if !ml.isConnected() {
		ml.connect()
	}
	// Hello
	err = ml.client.Hello(ml.hostname)
	if err != nil {
		return err
	}
	// Mail From
	err = ml.client.Mail(ml.mail.from.address)
	if err != nil {
		return err
	}
	// Rcpt To
	for _, addrs := range [][]*Address{ml.mail.to, ml.mail.cc, ml.mail.bcc} {
		for _, a := range addrs {
			err = ml.client.Rcpt(a.address)
			if err != nil {
				return err
			}
		}
	}
	// Data
	wc, err := ml.client.Data()
	if err != nil {
		return err
	}
	// fmt.Println(ml.mail.create())
	_, err = wc.Write([]byte(ml.mail.String()))
	if err != nil {
		return err
	}
	wc.Close()

	return nil
}

func (ml *MaiLetter) Reset() error {
	err := ml.client.Reset()
	if err != nil {
		return err
	}
	// ml.mail = NewMail()
	return nil
}

func (ml *MaiLetter) Noop() error {
	return ml.client.Noop()
}

func (ml *MaiLetter) Quit() error {
	return ml.client.Quit()
}

func (ml *MaiLetter) Close() error {
	return ml.client.Close()
}

func (ml *MaiLetter) isConnected() bool {
	if ml.client != nil {
		return true
	} else {
		return false
	}
}

func (ml *MaiLetter) connect() error {
	client, err := smtp.Dial(ml.dsn.Socket())
	if err != nil {
		return err
	}
	ml.client = client

	return nil
}
