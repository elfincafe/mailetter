package mailetter

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestNewAddress(t *testing.T) {
	cases := []struct {
		name    string
		address string
		typ     string
		err     error
	}{
		{"John Smith", "john@example.com", "*mailetter.Address", nil},
		{"Jane Smith", "jane@example.com", "*mailetter.Address", nil},
		{"", "jane@example.com", "*mailetter.Address", nil},
		{"", "janeatexample.com", "", errors.New("")},
	}

	for k, v := range cases {
		a := NewAddress(v.address, v.name)
		s := strings.Builder{}
		s.WriteString(fmt.Sprintf("[Case%d] ", k))
		length := len(s.String())
		if a.addr != v.address {
			msg := fmt.Sprintf("Addr: %s (%s), ", a.addr, v.address)
			s.WriteString(msg)
		}
		if a.name != v.name {
			msg := fmt.Sprintf("Name: %s (%s), ", a.name, v.name)
			s.WriteString(msg)
		}
		typ := reflect.TypeOf(a).String()
		if typ != v.typ {
			msg := fmt.Sprintf("Type: %s (%s), ", typ, v.typ)
			s.WriteString(msg)
		}
		if len(s.String()) > length {
			t.Error(strings.TrimRight(s.String(), ", "))
		}
	}
}

func TestAddress(t *testing.T) {
	type tcase struct {
		name    string
		address string
		err     error
	}
	cases := []tcase{}
	cases = append(cases, tcase{"John Smith", "john@example.com", nil})
	cases = append(cases, tcase{"Jane Smith", "jane@example.com", nil})
	cases = append(cases, tcase{"", "", errors.New("")})
	for k, v := range cases {
		a := NewAddress(v.address, v.name)
		s := strings.Builder{}
		s.WriteString(fmt.Sprintf("[Case%d]", k))
		length := len(s.String())
		if a.addr != v.address {
			s.WriteString(fmt.Sprintf("Addr: %s (%s), ", a.addr, v.address))
		}
		if len(s.String()) > length {
			t.Error(strings.Trim(s.String(), ", "))
		}
	}
}

func TestName(t *testing.T) {
	type tcase struct {
		name    string
		address string
		err     error
	}
	cases := []tcase{}
	cases = append(cases, tcase{"John Smith", "john@example.com", nil})
	cases = append(cases, tcase{"Jane Smith", "jane@example.com", nil})
	cases = append(cases, tcase{"", "", errors.New("")})
	for k, v := range cases {
		a := NewAddress(v.address, v.name)
		s := strings.Builder{}
		s.WriteString(fmt.Sprintf("[Case%d] ", k))
		if a.name != v.name {
			t.Errorf("[Case%d] %s != %s", k, a.name, v.name)
		}
	}
}

func TestAddrAngle(t *testing.T) {
	type tcase struct {
		address string
		expect  string
	}
	cases := []tcase{}
	cases = append(cases, tcase{"john@example.com", "<john@example.com>"})
	cases = append(cases, tcase{"jane@example.com", "<jane@example.com>"})
	cases = append(cases, tcase{"taro@example.com", "<taro@example.com>"})
	cases = append(cases, tcase{"", "<>"})
	for k, v := range cases {
		a := NewAddress(v.address, "")
		if a.String() != v.expect {
			t.Errorf("[Case%d] %s != %s", k, a.Angle(), v.expect)
		}
	}

}

func TestAddrString(t *testing.T) {
	type tcase struct {
		name string
		addr string
		mime string
	}
	cases := []tcase{}
	cases = append(cases, tcase{"John Smith", "john@example.com", "John Smith <john@example.com>"})
	cases = append(cases, tcase{"", "jane@example.com", "<jane@example.com>"})
	cases = append(cases, tcase{"山田 太郎", "taro@example.com", "=?UTF-8?B?5bGx55SwIOWkqumDjg==?= <taro@example.com>"})
	cases = append(cases, tcase{"", "", "<>"})
	for k, v := range cases {
		a := NewAddress(v.addr, v.name)
		if a.String() != v.mime {
			t.Errorf("[Case%d] %s != %s", k, a.String(), v.mime)
		}
	}
}
