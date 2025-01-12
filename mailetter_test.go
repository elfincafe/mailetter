package mailetter

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"
)

func TestEencode(t *testing.T) {
	cases := []struct {
		param   []byte
		flg     bool
		exptect []byte
	}{
		{
			[]byte("Mike Davice"),
			false,
			[]byte("Mike Davice"),
		},
		{
			[]byte("山田 太郎"),
			false,
			[]byte("5bGx55SwIOWkqumDjg=="),
		},
		{
			[]byte("Mike Davice"),
			true,
			[]byte("Mike Davice"),
		},
		{
			[]byte("山田 太郎"),
			true,
			[]byte("=?UTF-8?B?5bGx55SwIOWkqumDjg==?="),
		},
	}
	for i, c := range cases {
		result := encodeMime(c.param, c.flg)
		if !bytes.Equal(c.exptect, result) {
			t.Errorf(`[Case%d] Expect: %s,  Result: %s`, i+1, c.exptect, result)
		}
	}
}

func TestEncodeString(t *testing.T) {
	cases := []struct {
		param   string
		flg     bool
		exptect string
	}{
		{
			"Mike Davice",
			false,
			"Mike Davice",
		},
		{
			"山田 太郎",
			false,
			"5bGx55SwIOWkqumDjg==",
		},
		{
			"Mike Davice",
			true,
			"Mike Davice",
		},
		{
			"山田 太郎",
			true,
			"=?UTF-8?B?5bGx55SwIOWkqumDjg==?=",
		},
	}

	for i, c := range cases {
		result := encodeMimeString(c.param, c.flg)
		if c.exptect != result {
			t.Errorf(`[Case%d] Expect: %s,  Result: %s`, i+1, c.exptect, result)
		}
	}
}

func TestBorder(t *testing.T) {
	length := 24
	border := border(length)
	re := regexp.MustCompile(fmt.Sprintf(`-{12}[0-9a-zA-Z]{%d}`, length))
	if !re.MatchString(border) {
		t.Errorf(`Invalid Border %s`, border)
	}
}
