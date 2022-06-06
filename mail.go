package mailetter

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"
	"bytes"
)

type Mail struct {
	orders      []string
	headers     map[string]map[string]string
	from        *Addr
	returnPath	*Addr
	to          []*Addr
	cc          []*Addr
	bcc         []*Addr
	subj        string
	body        string
	html        string
	attachments []*Attachment
	vars        map[string]interface{}
}

func NewMail() *Mail {
	m := new(Mail)
	m.orders = []string{}
	m.headers = map[string]map[string]string{}
	m.from = nil
	m.returnPath = nil
	m.to = []*Addr{}
	m.cc = []*Addr{}
	m.bcc = []*Addr{}
	m.subj = ""
	m.body = ""
	m.html = ""
	m.attachments = []*Attachment{}
	m.vars = map[string]interface{}{}
	return m
}

func (m *Mail) Header(key, val string) {
	k := strings.ToLower(key)
	if _, ok := m.headers[k]; !ok {
		m.orders = append(m.orders, k)
	}
	m.headers[k] = map[string]string{"key": key, "val": val}
}

func (m *Mail) To(addr *Addr) {
	m.to = append(m.to, addr)
}

func (m *Mail) Cc(addr *Addr) {
	m.cc = append(m.cc, addr)
}

func (m *Mail) Bcc(addr *Addr) {
	m.bcc = append(m.bcc, addr)
}

func (m *Mail) From(addr *Addr) {
	m.from = addr
	if m.returnPath == nil {
		m.returnPath = addr
	}
}

func (m *Mail) Subject(subj string) {
	m.subj = subj
}

func (m *Mail) Set(key string, val interface{}) {
	m.vars[key] = val
}

func (m *Mail) Body(body string) {
	m.body = body
}

func (m *Mail) LoadBody(path string) error {
	subj, body, err := m.load(path)
	if err != nil {
		return err
	}
	m.Subject(string(subj))
	m.Body(string(body))
	return nil
}

func (m *Mail) Html(html string) {
	m.html = html
}

func (m *Mail) LoadHtml(path string) error {
	subj, html, err := m.load(path)
	if err != nil {
		return err
	}
	m.Subject(string(subj))
	m.Html(string(html))
	return nil
}

func (m *Mail) Attach(a *Attachment) {
	m.attachments = append(m.attachments, a)
}

func (m *Mail) load(path string) ([]byte, []byte, error) {
	if _, err := os.Lstat(path); os.IsNotExist(err) {
		return []byte{}, []byte{}, err
	}
	tpl := template.Must(template.ParseFiles(path))
	buf := bytes.NewBuffer(make([]byte, 0, 10240))
	err := tpl.Execute(buf, m.vars)
	if err != nil {
		return []byte{}, []byte{}, err
	}
	content := buf.Bytes()
	idx := 0
	// Subject
	pos := bytes.Index(content, []byte("\n"))
	if idx<0 {
		return []byte{}, []byte{}, err
	}
	idx += pos
	subj := bytes.Trim(content[:idx], "\r\n")
	// Body
	pos = bytes.Index(content[idx+1:], []byte("\n"))
	if idx<0 { // Skip Empty line
		return []byte{}, []byte{}, err
	}
	idx += pos
	body := content[idx+1:]

	return subj, body, nil
}

func (m *Mail) create() string {
	content := strings.Builder{}
	header := ""
	// Header
	for _, k := range m.orders {
		hdr := m.headers[k]
		content.WriteString(fmt.Sprintf("%s: %s\r\n", hdr["key"], hdr["val"]))
	}
	// Content-Type
	header = fmt.Sprintf("Content-Type: %s\r\n", "text/plain; charset=UTF-8")
	content.WriteString(header)
	// Date
	header = fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	content.WriteString(header)
	// From
	if len(m.from.name) > 0 {
		header = fmt.Sprintf("From: %s <%s>\r\n", encode(m.from.name), m.from.addr)
	} else {
		header = fmt.Sprintf("From: <%s>\r\n", m.from.addr)
	}
	content.WriteString(header)
	// To
	to := strings.Builder{}
	for k, v := range m.to {
		indent := ""
		if k > 0 {
			indent = "    "
		}
		if len(v.name) > 0 {
			to.WriteString(fmt.Sprintf("%s%s <%s>,\r\n", indent, encode(v.name), v.addr))
		} else {
			to.WriteString(fmt.Sprintf("%s<%s>,\r\n", indent, v.addr))
		}
	}
	header = fmt.Sprintf("To: %s\r\n", strings.TrimRight(to.String(), "\r\n,"))
	content.WriteString(header)
	// Cc
	cc := strings.Builder{}
	for k, v := range m.cc {
		indent := ""
		if k > 0 {
			indent = "    "
		}
		if len(v.name) > 0 {
			cc.WriteString(fmt.Sprintf("%s%s <%s>,\r\n", indent, encode(v.name), v.addr))
		} else {
			cc.WriteString(fmt.Sprintf("%s<%s>,\r\n", indent, v.addr))
		}
	}
	header = fmt.Sprintf("Cc: %s\r\n", strings.TrimRight(cc.String(), "\r\n,"))
	if len(m.cc) > 0 {
		content.WriteString(header)
	}
	// Subject
	content.WriteString(fmt.Sprintf("Subject: %s\r\n", encode(m.subj)))
	content.WriteString("\r\n")
	// Body
	content.WriteString(m.body)

	return content.String()
}
