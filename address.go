package mailetter

import (
	"fmt"
	"net/mail"
	"strings"
)

type Address struct {
	addr string
	name string
	err  error
}

func NewAddress(address, name string) *Address {
	address = strings.Trim(address, white_space)
	name = strings.Trim(name, white_space)
	a := &Address{addr: "", name: "", err: nil}
	pa, err := mail.ParseAddress(fmt.Sprintf("%s <%s>", name, address))
	if err != nil {
		a.err = err
	} else {
		a.addr = pa.Address
		a.name = pa.Name
	}
	return a
}

func (a *Address) Address() string {
	return a.addr
}

func (a *Address) Name() string {
	return a.name
}

func (a *Address) Error() error {
	return a.err
}

func (a *Address) Angle() string {
	return fmt.Sprintf("<%s>", a.addr)
}

func (a *Address) String() string {
	if len(a.name) > 0 {
		return fmt.Sprintf("%s %s", EncodeMimeString(a.name, true), a.Angle())
	} else {
		return a.Angle()
	}
}
