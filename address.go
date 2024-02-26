package mailetter

import (
	"fmt"
	"net/mail"
	"strings"
)

type Address struct {
	address string
	name    string
}

func NewAddress(address, name string) (*Address, error) {
	address = strings.Trim(address, " \r\n\t\v")
	name = strings.Trim(name, " \r\n\t\v")
	pa, err := mail.ParseAddress(fmt.Sprintf("%s <%s>", name, address))
	if err != nil {
		return nil, err
	}
	a := new(Address)
	a.address = pa.Address
	a.name = pa.Name
	return a, nil
}

func (a *Address) Address() string {
	return a.address
}

func (a *Address) Name() string {
	return a.name
}

func (a *Address) Angle() string {
	return fmt.Sprintf("<%s>", a.address)
}

func (a *Address) String() string {
	if len(a.name) > 0 {
		return fmt.Sprintf("%s %s", EncodeMimeString(a.name, true), a.Angle())
	} else {
		return a.Angle()
	}
}
