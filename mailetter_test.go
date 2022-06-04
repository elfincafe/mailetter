package mailetter

import (
	"fmt"
	"strings"
	"testing"
)

const (
	DSN       = "smtp://smtp.example.com:25"
	FROM_ADDR = "test@example.com"
	FROM_NAME = ""
)

func TestIsConnected(t *testing.T) {
	goodDsn := DSN
	opts := map[string]interface{}{}
	m, _ := New(goodDsn, opts)
	m.connect()
	// fmt.Println(m.client)
	if !m.isConnected() {
		t.Errorf("Connection Failed. Something is wrong.")
	}
	badDsn := strings.ReplaceAll(DSN, ":25", ":12345")
	m, _ = New(badDsn, opts)
	m.connect()
	if m.isConnected() {
		t.Errorf("Connection Success. Something is wrong.")
	}
}

func TestConnect(t *testing.T) {
	opts := map[string]interface{}{}
	_, err := New(DSN, opts)
	if err != nil {
		fmt.Println(err)
	}
}
