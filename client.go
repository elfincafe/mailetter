package mailetter

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

type Client struct {
	dsn       *dsn
	conn      *smtp.Client
	localName string
	auth      smtp.Auth
	data      *data
}

func New(dsnStr string) *Client {
	c := new(Client)
	c.dsn = newDsn(dsnStr)
	c.conn = nil
	c.localName = "localhost"
	c.auth = nil
	data := newData()
	data.reset()
	c.data = data
	return c
}

func (c *Client) LocalName(localName string) {
	c.localName = localName
}

func (c *Client) AuthByPlain(username, password string) {
	if c.dsn.host == "" {
		_ = c.dsn.parse()
	}
	c.auth = smtp.PlainAuth(username, username, password, c.dsn.host)
}

func (c *Client) AuthByLogin(username, password string) {
	auth := new(Login)
	auth.username = username
	auth.password = password
	c.auth = auth
}

func (c *Client) AuthByCramMd5(username, secret string) {
	c.auth = smtp.CRAMMD5Auth(username, secret)
}

func (c *Client) Header(key, value string) error {
	return c.data.setHeader(key, value)
}

func (c *Client) From(addr, name string) error {
	return c.data.setFrom(newAddress(addr, name))
}

func (c *Client) To(addr, name string) error {
	return c.data.setTo(newAddress(addr, name))
}

func (c *Client) Cc(addr, name string) error {
	return c.data.setCc(newAddress(addr, name))
}

func (c *Client) Bcc(addr, name string) error {
	return c.data.setBcc(newAddress(addr, name))
}

func (c *Client) Subject(subject string) {
	c.data.setSubject(subject)
}

func (c *Client) Body(body string) {
	c.data.setBody(strings.NewReader(body))
}

func (c *Client) Set(key, value string) {
	c.data.setValue(key, value)
}

func (c *Client) Send() error {

	var err error
	if err = c.connect(); err != nil {
		return err
	}

	// Mail From
	if c.data.from == nil {
		return fmt.Errorf(`a from address is required`)
	}
	err = c.conn.Mail(c.data.from.addr)
	if err != nil {
		return err
	}

	// Rcpt To
	if len(c.data.to) == 0 {
		return fmt.Errorf(`at least one recipient is required`)
	}
	for _, addrs := range [][]*Address{c.data.to, c.data.cc, c.data.bcc} {
		for _, a := range addrs {
			fmt.Println(a.Angle())
			err = c.conn.Rcpt(a.addr)
			if err != nil {
				return err
			}
		}
	}

	// Data
	wc, err := c.conn.Data()
	if err != nil {
		return err
	}
	fmt.Println(c.data.String())
	_, err = wc.Write([]byte(c.data.String()))
	if err != nil {
		return err
	}
	wc.Close()

	return nil
}

func (c *Client) Reset() error {
	c.data.reset()
	return c.conn.Reset()
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
	var err error
	if c.isConnected() {
		return nil
	}
	err = c.dsn.parse()
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
	if err != nil {
		fmt.Println("[1]", err)
		return err
	}
	// Hello
	err = c.conn.Hello(c.localName)
	if err != nil {
		fmt.Println("[2]", err)
		return err
	}
	// Auth
	fmt.Println(c.dsn)
	if c.auth != nil {
		err = c.conn.Auth(c.auth)
		if err != nil {
			fmt.Println("[3]", err)
			return err
		}
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
