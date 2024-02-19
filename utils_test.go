package mailetter

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"
)

func TestEencode(t *testing.T) {
	type case1 struct {
		param   []byte
		flg     bool
		exptect []byte
	}
	cases1 := []case1{}
	cases1 = append(cases1, case1{[]byte("Mike Davice"), false, []byte("Mike Davice")})
	cases1 = append(cases1, case1{[]byte("山田 太郎"), false, []byte("5bGx55SwIOWkqumDjg==")})
	cases1 = append(cases1, case1{[]byte("Mike Davice"), true, []byte("Mike Davice")})
	cases1 = append(cases1, case1{[]byte("山田 太郎"), true, []byte("=?UTF-8?B?5bGx55SwIOWkqumDjg==?=")})

	for i, c := range cases1 {
		result := EncodeMime(c.param, c.flg)
		if !bytes.Equal(c.exptect, result) {
			t.Errorf(`[Case%d] Expect: %s,  Result: %s`, i+1, c.exptect, result)
		}
	}
}

func TestEncodeString(t *testing.T) {
	type case2 struct {
		param   string
		flg     bool
		exptect string
	}
	cases2 := []case2{}
	cases2 = append(cases2, case2{"Mike Davice", false, "Mike Davice"})
	cases2 = append(cases2, case2{"山田 太郎", false, "5bGx55SwIOWkqumDjg=="})
	cases2 = append(cases2, case2{"Mike Davice", true, "Mike Davice"})
	cases2 = append(cases2, case2{"山田 太郎", true, "=?UTF-8?B?5bGx55SwIOWkqumDjg==?="})

	for i, c := range cases2 {
		result := EncodeMimeString(c.param, c.flg)
		if c.exptect != result {
			t.Errorf(`[Case%d] Expect: %s,  Result: %s`, i+1, c.exptect, result)
		}
	}
}

func TestBorder(t *testing.T) {
	length := 24
	border := Border(length)
	re := regexp.MustCompile(fmt.Sprintf(`-{12}[0-9a-zA-Z]{%d}`, length))
	if !re.MatchString(border) {
		t.Errorf(`Invalid Border %s`, border)
	}
}
