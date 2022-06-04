package mailetter

import (
	"encoding/base64"
	"strings"
	"math/rand"
	"fmt"
	"time"
)
func encode(s string) string {
	needsEnc := false
	for _, v := range s {
		if v > 127 {
			needsEnc = true
			break
		}
	}
	if !needsEnc {
		return s
	}
	return fmt.Sprintf("=?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(s)))
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
	sb.WriteString(strings.Repeat("-", 12))
	for i := 0; i < length; i++ {
		idx := rand.Intn(l - 1)
		sb.WriteString(s[idx])
	}
	return sb.String()
}
