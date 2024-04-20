package mailetter

import (
	"crypto/tls"
	"net/smtp"
)

const (
	br          = "\r\n"
	white_space = " \r\n\t\v\b"
	shouldBr    = 78
	mustBr      = 998
)

type MaiLetter struct {
	dsn       *dsn
	client    *smtp.Client
	localName string
	tlsConfig *tls.Config
}

func New(dsn string) (*MaiLetter, error) {
	oDsn, err := newDsn(dsn)
	if err != nil {
		return nil, err
	}
	ml := new(MaiLetter)
	ml.dsn = oDsn
	ml.client = nil
	ml.localName = "localhost"
	ml.tlsConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         oDsn.Host(),
	}
	return ml, nil
}

func (ml *MaiLetter) LocalName(localName string) {
	ml.localName = localName
}

func (ml *MaiLetter) Send(m *Mail) error {

	var err error
	if err = ml.connect(); err != nil {
		return err
	}

	// Hello
	err = ml.client.Hello(ml.localName)
	if err != nil {
		return err
	}

	// Mail From
	err = ml.client.Mail(m.from.addr)
	if err != nil {
		return err
	}

	// Rcpt To
	for _, addrs := range [][]*Addr{m.to, m.cc, m.bcc} {
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
	if ml.client != nil {
		return ml.client.Close()
	}
	return nil
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
	var err error
	if ml.dsn.IsSsl() {
		err = ml.connectWithSsl()
	} else {
		err = ml.connectWithoutSsl()
	}
	return err
}

func (ml *MaiLetter) connectWithoutSsl() error {
	client, err := smtp.Dial(ml.dsn.Socket())
	if err != nil {
		return err
	}
	ml.client = client
	return nil
}

func (ml *MaiLetter) connectWithSsl() error {
	conn, err := tls.Dial("tcp", ml.dsn.Socket(), ml.tlsConfig)
	if err != nil {
		return err
	}
	client, err := smtp.NewClient(conn, ml.dsn.Host())
	if err != nil {
		return err
	}
	ml.client = client
	return nil
}

func (ml *MaiLetter) connectAndStartTls() error {
	var err error
	err = ml.connectWithoutSsl()
	if err != nil {
		return err
	}
	err = ml.client.StartTLS(ml.tlsConfig)
	if err != nil {
		return err
	}
	return nil
}
