package mailetter

import (
	"net/smtp"
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
}

func New(dsn *Dsn) *MaiLetter {
	ml := new(MaiLetter)
	ml.dsn = dsn
	ml.client = nil
	ml.hostname = ""
	return ml
}

func (ml *MaiLetter) Hostname(hostname string) {
	ml.hostname = hostname
}

func (ml *MaiLetter) Send(m *Mail) error {

	var err error
	if err = ml.connect(); err != nil {
		return err
	}

	// Hello
	err = ml.client.Hello(ml.hostname)
	if err != nil {
		return err
	}
	// Mail From
	err = ml.client.Mail(m.from.addr)
	if err != nil {
		return err
	}
	// Rcpt To
	for _, addrs := range [][]*Address{m.to, m.cc, m.bcc} {
		for _, a := range addrs {
			err = ml.client.Rcpt(a.addr)
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
	_, err = wc.Write([]byte(m.String()))
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
	return nil
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
	if ml.isConnected() {
		return nil
	}

	client, err := smtp.Dial(ml.dsn.Socket())
	if err != nil {
		return err
	}
	ml.client = client

	return nil
}
