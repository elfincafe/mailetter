package mailetter

import (
	"testing"
	"regexp"
)

func TestBorder(t *testing.T) {
	border := border()
	re := regexp.MustCompile(`-{12}[0-9a-zA-Z]{24}`)
	if !re.MatchString(border) {
		t.Errorf(`Invalid Border %s`, border)
	}
}
