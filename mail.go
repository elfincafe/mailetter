package mailetter

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"
	"time"
)

type header struct {
	key string
	val string
}

type Mail struct {
	headers    map[string]header
	from       *Addr
	returnPath *Addr
	replyTo    *Addr
	to         []*Addr
	cc         []*Addr
	bcc        []*Addr
	subject    *template.Template
	body       *template.Template
	vars       map[string]any
}

func NewMail(from *Addr) *Mail {
	m := new(Mail)
	m.from = from
	m.returnPath = from
	m.replyTo = from
	m.Reset()
	return m
}

func (m *Mail) Reset() {
	m.headers = map[string]header{}
	m.to = []*Addr{}
	m.cc = []*Addr{}
	m.bcc = []*Addr{}
	m.subject = nil
	m.body = nil
	m.vars = map[string]any{}
}

func (m *Mail) Header(key, val string) {
	key = strings.Trim(key, white_space)
	val = strings.Trim(val, white_space)
	for old, new := range map[string]string{"\r": "", "\n": ""} {
		key = strings.ReplaceAll(key, old, new)
		val = strings.ReplaceAll(val, old, new)
	}
	excepts := map[string]bool{
		"content-type": true,
		"date":         true,
		"from":         true,
		"reply-to":     true,
		"return-path":  true,
		"to":           true,
		"cc":           true,
		"bcc":          true,
		"subject":      true,
	}
	lowerKey := strings.ToLower(key)
	if _, ok := excepts[lowerKey]; ok {
		return
	}

	m.headers[lowerKey] = header{key: key, val: val}
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

func (m *Mail) ReturnPath(addr *Addr) {
	m.returnPath = addr
}

func (m *Mail) ReplyTo(addr *Addr) {
	m.replyTo = addr
}

func (m *Mail) Subject(subject string) {
	for old, new := range map[string]string{"\r": "", "\n": ""} {
		subject = strings.ReplaceAll(subject, old, new)
	}
	m.subject = template.Must(template.New("Subject").Parse(subject))
}

func (m *Mail) Body(r io.Reader) {
	body, _ := io.ReadAll(r)
	rplr := strings.NewReplacer("\r\n", "\n", "\r", "\n", "\n", "\r\n")
	bodyText := rplr.Replace(string(body))
	m.body = template.Must(template.New("Body").Parse(bodyText))
}

func (m *Mail) Set(key string, val any) {
	m.vars[key] = val
}

func (m *Mail) String() string {
	sb := strings.Builder{}
	line := strings.Builder{}
	// buf := bytes.NewBuffer(make([]byte, 10240))

	// Headers
	line.Reset()
	for _, v := range m.headers {
		line.WriteString(fmt.Sprintf("%s: %s\r\n", v.key, encodeMimeString(v.val, true)))
		sb.WriteString(line.String())
	}

	// Content-Type
	sb.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	// Date
	sb.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	// From
	sb.WriteString(fmt.Sprintf("From: %s\r\n", encodeMimeString(m.from.String(), true)))
	// To, Cc
	rcpts := map[string][]*Addr{"To": m.to, "Cc": m.cc}
	for label, addrs := range rcpts {
		if len(addrs) == 0 {
			continue
		}
		line.Reset()
		line.WriteString(label)
		line.WriteString(":")
		for k, v := range addrs {
			if k == 0 {
				line.WriteString(" ")
			} else {
				line.WriteString("    ")
			}
			line.WriteString(v.String())
			if len(addrs) > k+1 {
				line.WriteString(",")
			}
			line.WriteString("\r\n")
		}
		sb.WriteString(line.String())
	}

	// Subject
	if m.subject == nil {
		m.Subject("")
	}
	buf := bytes.NewBuffer([]byte{})
	m.subject.Execute(buf, m.vars)
	subject := encodeMimeString(buf.String(), true)
	line.Reset()
	line.WriteString("Subject: ")
	line.WriteString(encodeMimeString(subject, true))
	line.WriteString("\r\n")
	sb.WriteString(line.String())

	// Body
	buf.Reset()
	m.body.Execute(buf, m.vars)
	line.Reset()
	line.WriteString("\r\n")
	line.WriteString(buf.String())
	sb.WriteString(line.String())

	return sb.String()
}
