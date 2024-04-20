package mailetter

import (
	"fmt"
	"net/mail"
	"strings"
)

type Addr struct {
	addr string
	name string
	err  error
}

func NewAddr(address, name string) *Addr {
	address = strings.Trim(address, white_space)
	name = strings.Trim(name, white_space)
	a := &Addr{addr: "", name: "", err: nil}
	pa, err := mail.ParseAddress(fmt.Sprintf("%s <%s>", name, address))
	if err != nil {
		a.err = err
	} else {
		a.addr = pa.Address
		a.name = pa.Name
	}
	return a
}

func (a *Addr) Address() string {
	return a.addr
}

func (a *Addr) Name() string {
	return a.name
}

func (a *Addr) Error() error {
	return a.err
}

func (a *Addr) Angle() string {
	return fmt.Sprintf("<%s>", a.addr)
}

func (a *Addr) String() string {
	if len(a.name) > 0 {
		return fmt.Sprintf("%s %s", encodeMimeString(a.name, true), a.Angle())
	} else {
		return a.Angle()
	}
}
