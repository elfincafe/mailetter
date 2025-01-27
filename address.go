package mailetter

import (
	"fmt"
	"net/mail"
	"strings"
)

type Address struct {
	addr string
	name string
}

func NewAddress(address, name string) *Address {
	address = strings.Trim(address, white_space)
	name = strings.Trim(name, white_space)
	a := new(Address)
	a.addr = address
	a.name = name
	return a
}

func (a *Address) Angle() string {
	return fmt.Sprintf("<%s>", a.addr)
}

func (a *Address) String() string {
	if len(a.name) > 0 {
		return fmt.Sprintf("%s %s", encodeMimeString(a.name, true), a.Angle())
	} else {
		return a.Angle()
	}
}

func (a *Address) parse() error {
	_, err := mail.ParseAddress(a.Angle())
	if err != nil {
		return err
	}
	return nil
}
