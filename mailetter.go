package mailetter

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"
	"time"
)

type MaiLetter struct {
	dsn    string
	client string
	to     [][2]string
	cc     [][2]string
	bcc    []string
	subj   string
	body   string
	from   [2]string
	atts   []*Att
	auth   Auth
	border string
}

type Att struct {
	content []byte
	name    string
	path    string
}

type Auth struct {
	user     string
	password string
}

func New(dsn string, addr string, name string) *MaiLetter {
	m := new(MaiLetter)
	m.dsn = dsn
	m.from = [2]string{addr, name}
	m.border = fmt.Sprintf("----------%s", border())
	return m
}

func NewAttFrom(path string, name string) *Att {
	a := new(Att)
	a.path = path
	if len(name) > 0 {
		a.name = name
	} else {
		a.name = filepath.Base(path)
	}
	return a
}

func NewAtt(name string, content []byte) *Att {
	a := new(Att)
	a.content = content
	a.name = name
	return a
}

func (m *MaiLetter) Authenticate(auth Auth) {
	m.auth = auth
}

func (m *MaiLetter) To(addr string, name string) {
	to := [2]string{addr, name}
	m.to = append(m.to, to)
}

func (m *MaiLetter) Cc(addr string, name string) {
	cc := [2]string{addr, name}
	m.cc = append(m.cc, cc)
}

func (m *MaiLetter) Bcc(addr string) {
	m.bcc = append(m.bcc, addr)
}

func (m *MaiLetter) Subj(subj string) {
	m.subj = subj
}

func (m *MaiLetter) Body(body string) {
	m.body = body
}

func (m *MaiLetter) Attach(a *Att) {
	m.atts = append(m.atts, a)
}

func (m *MaiLetter) Send() bool {
	return true
}

func (m *MaiLetter) Reset() bool {
	return true
}

func (m *MaiLetter) Noop() bool {
	return true
}

func encode(s string) string {
	should := false
	for {

	}
	if !should {
		return s
	}
	return base64.StdEncoding.EncodeToString([]byte(s))
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
	for i := 0; i < length; i++ {
		idx := rand.Intn(l - 1)
		sb.WriteString(s[idx])
	}
	return strings.Repeat("-", 6) + sb.String()
}
