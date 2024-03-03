package mailetter

import (
	"bytes"
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
	from       *Address
	returnPath *Address
	replyTo    *Address
	to         []*Address
	cc         []*Address
	bcc        []*Address
	subject    *template.Template
	body       *template.Template
	vars       map[string]any
}

func NewMail(from *Address) *Mail {
	m := new(Mail)
	m.from = from
	m.returnPath = from
	m.replyTo = from
	m.Reset()
	return m
}

func (m *Mail) Reset() {
	m.headers = map[string]header{}
	m.to = []*Address{}
	m.cc = []*Address{}
	m.bcc = []*Address{}
	m.subject = nil
	m.body = nil
	m.vars = map[string]any{}
}

func (m *Mail) Header(key, val string) {
	key = strings.Trim(key, white_space)
	val = strings.Trim(val, white_space)
	m.headers[strings.ToLower(key)] = header{key: key, val: val}
}

func (m *Mail) To(addr *Address) {
	m.to = append(m.to, addr)
}

func (m *Mail) Cc(addr *Address) {
	m.cc = append(m.cc, addr)
}

func (m *Mail) Bcc(addr *Address) {
	m.bcc = append(m.bcc, addr)
}

func (m *Mail) ReturnPath(addr *Address) {
	m.returnPath = addr
}

func (m *Mail) ReplyTo(addr *Address) {
	m.replyTo = addr
}

func (m *Mail) Subject(subject string) {
	m.subject = template.Must(template.New("Subject").Parse(subject))
}

func (m *Mail) Body(body string) {
	m.body = template.Must(template.New("Body").Parse(body))
}

func (m *Mail) Set(key string, val any) {
	m.vars[key] = val
}

func (m *Mail) String() string {
	sb := strings.Builder{}
	line := strings.Builder{}
	buf := bytes.NewBuffer(make([]byte, 10240))

	// Content-Type
	m.headers["content-type"] = header{"Content-Type", "text/plain; charset=UTF-8"}
	// Date
	m.headers["date"] = header{"Date", time.Now().Format(time.RFC1123Z)}
	// From
	m.headers["from"] = header{"From", m.from.String()}
	// To
	line.Reset()
	line.WriteString("To:")
	for k, v := range m.to {
		angle := v.String()
		if line.Len()+len(angle)+2 > should_br {
			sb.WriteString(line.String() + br)
			line.Reset()
			line.WriteString(" ")
		}
		line.WriteString(" " + angle + ",")
		if len(m.to) >= k-1 {
			sb.WriteString(line.String())
		}
	}
	m.headers["to"] = header{"To", strings.TrimRight(sb.String(), ","+white_space)}

	// Cc
	line.Reset()
	line.WriteString("Cc:")
	for k, v := range m.to {
		angle := v.String()
		if line.Len()+len(angle)+2 > should_br {
			sb.WriteString(line.String() + br)
			line.Reset()
			line.WriteString(" ")
		}
		line.WriteString(" " + angle + ",")
		if len(m.to) >= k-1 {
			sb.WriteString(line.String())
		}
	}
	m.headers["cc"] = header{"Cc", strings.TrimRight(sb.String(), ","+white_space)}

	// Subject
	m.subject.Execute(buf, m.vars)
	subject := EncodeMimeString(buf.String(), true)
	line.Reset()
	line.WriteString("Subject: ")
	m.headers["subject"] = header{"Subject", subject}

	// Body
	// body := bytes.NewBuffer()
	// content.WriteString(m.body)

	// return content.String()
	return ""
}
