package mailetter

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestNewAddr(t *testing.T) {
	type tcase struct {
		name string
		addr string
		typ  string
	}
	cases := []tcase{}
	cases = append(cases, tcase{"John Smith", "john@example.com", "*mailetter.Addr"})
	cases = append(cases, tcase{"Jane Smith", "jane@example.com", "*mailetter.Addr"})
	cases = append(cases, tcase{"", "jane@example.com", "*mailetter.Addr"})
	for k, v := range cases {
		a := NewAddr(v.addr, v.name)
		e := false
		s := strings.Builder{}
		s.WriteString(fmt.Sprintf("[Case%d] ", k))
		if a.addr != v.addr {
			e = true
			msg := fmt.Sprintf("Addr: %s (%s), ", a.addr, v.addr)
			s.WriteString(msg)
		}
		if a.name != v.name {
			e = true
			msg := fmt.Sprintf("Name: %s (%s), ", a.name, v.name)
			s.WriteString(msg)
		}
		typ := reflect.TypeOf(a).String()
		if typ != v.typ {
			e = true
			msg := fmt.Sprintf("Type: %s (%s), ", typ, v.typ)
			s.WriteString(msg)
		}
		if e {
			t.Error(strings.TrimRight(s.String(), ", "))
		}
	}
}

func TestAddr(t *testing.T) {
	type tcase struct {
		name string
		addr string
	}
	cases := []tcase{}
	cases = append(cases, tcase{"John Smith", "john@example.com"})
	cases = append(cases, tcase{"Jane Smith", "jane@example.com"})
	cases = append(cases, tcase{"", ""})
	for k, v := range cases {
		a := NewAddr(v.addr, v.name)
		if a.addr != v.addr {
			t.Errorf("[Case%d] %s != %s", k, a.addr, v.addr)
		}
	}
}

func TestName(t *testing.T) {
	type tcase struct {
		name string
		addr string
	}
	cases := []tcase{}
	cases = append(cases, tcase{"John Smith", "john@example.com"})
	cases = append(cases, tcase{"Jane Smith", "jane@example.com"})
	cases = append(cases, tcase{"", ""})
	for k, v := range cases {
		a := NewAddr(v.addr, v.name)
		if a.name != v.name {
			t.Errorf("[Case%d] %s != %s", k, a.name, v.name)
		}
	}
}

func TestAddrString(t *testing.T) {
	type tcase struct {
		name   string
		addr   string
		expect string
	}
	cases := []tcase{}
	cases = append(cases, tcase{"John Smith", "john@example.com", "John Smith <john@example.com>"})
	cases = append(cases, tcase{"", "jane@example.com", "<jane@example.com>"})
	cases = append(cases, tcase{"山田 太郎", "taro@example.com", "=?UTF-8?B?5bGx55SwIOWkqumDjg==?= <taro@example.com>"})
	cases = append(cases, tcase{"", "", "<>"})
	for k, v := range cases {
		a := NewAddr(v.addr, v.name)
		if a.String() != v.expect {
			t.Errorf("[Case%d] %s != %s", k, a.String(), v.expect)
		}
	}
}
