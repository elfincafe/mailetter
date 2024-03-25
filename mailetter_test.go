package mailetter

import (
	"reflect"
	"testing"
)

func TestMaiLetterNew(t *testing.T) {
	dsn, _ := NewDsn("smtp://example.com")
	cases := []struct {
		dsn      *Dsn
		expected string
	}{
		{dsn, "*mailetter.MaiLetter"},
	}
	for k, v := range cases {
		m := New(v.dsn)
		if reflect.TypeOf(m).String() != v.expected {
			t.Errorf(`[Case%d] %v`, k, reflect.TypeOf(m))
		}
	}
}

func TestMaiLetterLocalName(t *testing.T) {
	cases := []struct {
		call bool
		name string
	}{
		{false, "localhost.localdomain"},
		{true, "mail.example.com"},
	}
	for k, v := range cases {
		dsn, _ := NewDsn("smtp://localhost")
		m := New(dsn)
		if v.call {
			m.LocalName(v.name)
		}
		if m.localName != v.name {
			t.Errorf(`[Case%d] %s (%s)`, k, m.localName, v.name)
		}
	}
}
