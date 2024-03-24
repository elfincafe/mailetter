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
