package mailetter

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

const (
	br          = "\r\n"
	white_space = " \r\n\t\v\b"
	shouldBr    = 78
	mustBr      = 998
)

type Client struct {
	dsn       *dsn
	conn      *smtp.Client
	localName string
	auth      AuthInterface
	from      *Address
	to        []*Address
	cc        []*Address
	bcc       []*Address
	subj      string
	body      string
}

func New(dsnStr string) *Client {
	c := new(Client)
	c.dsn = newDsn(dsnStr)
	c.conn = nil
	c.localName = "localhost"
	c.auth = nil
	c.from = nil
	c.to = []*Address{}
	c.cc = []*Address{}
	c.bcc = []*Address{}
	c.subj = ""
	c.body = ""
	return c
}

func (c *Client) LocalName(localName string) {
	c.localName = localName
}

func (c *Client) Auth(auth AuthInterface) {
	c.auth = auth
}

func (c *Client) From(addr, name string) error {
	c.from = NewAddress(addr, name)
	if err := c.from.parse(); err != nil {
		return fmt.Errorf(`the addess for "From:" is invalid (%s)`, err.Error())
	}
	return nil
}

func (c *Client) To(addr, name string) error {
	to := NewAddress(addr, name)
	if err := to.parse(); err != nil {
		return fmt.Errorf(`the addess for "To:" is invalid (%s)`, err.Error())
	}
	c.to = append(c.to, to)
	return nil
}

func (c *Client) Cc(addr, name string) error {
	cc := NewAddress(addr, name)
	if err := cc.parse(); err != nil {
		return fmt.Errorf(`the addess for "Cc:" is invalid (%s)`, err.Error())
	}
	c.cc = append(c.cc, cc)
	return nil
}

func (c *Client) Bcc(addr, name string) error {
	bcc := NewAddress(addr, name)
	if err := bcc.parse(); err != nil {
		return fmt.Errorf(`the addess for "Bcc:" is invalid (%s)`, err.Error())
	}
	c.bcc = append(c.bcc, bcc)
	return nil
}

func (c *Client) Subject(subject string) {
	c.subj = subject
}

func (c *Client) Body(body string) {
	c.body = body
}

func (ml *Client) Send(m *Mail) error {

	var err error
	if err = ml.connect(); err != nil {
		return err
	}

	// Hello
	err = ml.conn.Hello(ml.localName)
	if err != nil {
		return err
	}

	// Mail From
	err = ml.conn.Mail(m.from.addr)
	if err != nil {
		return err
	}

	// Rcpt To
	for _, addrs := range [][]*Address{m.to, m.cc, m.bcc} {
		for _, a := range addrs {
			err = ml.conn.Rcpt(a.addr)
			if err != nil {
				return err
			}
		}
	}

	// Data
	wc, err := ml.conn.Data()
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

func (c *Client) Close() error {
	return c.conn.Quit()
}

func (c *Client) isConnected() bool {
	if c.conn != nil {
		return true
	} else {
		return false
	}
}

func (c *Client) connect() error {
	if c.isConnected() {
		return nil
	}
	err := c.dsn.parse()
	if err != nil {
		return err
	}
	switch c.dsn.scheme {
	case "smtps":
		c.conn, err = c.connectSmtps(c.dsn)
	case "smtp+tls":
		c.conn, err = c.connectWithTls(c.dsn)
	case "smtp":
		c.conn, err = c.connectSmtp(c.dsn)
	}

	return err
}

func (c *Client) connectSmtps(d *dsn) (*smtp.Client, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         d.host,
	}
	conn, err := tls.Dial("tcp", d.Socket(), tlsConfig)
	if err != nil {
		return nil, err
	}
	return smtp.NewClient(conn, d.host)
}

func (c *Client) connectSmtp(d *dsn) (*smtp.Client, error) {
	return smtp.Dial(d.Socket())
}

func (c *Client) connectWithTls(d *dsn) (*smtp.Client, error) {
	conn, err := c.connectSmtp(d)
	if err != nil {
		return nil, err
	}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         d.host,
	}
	err = conn.StartTLS(tlsConfig)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
