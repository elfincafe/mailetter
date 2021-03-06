package mailetter

import (
	"fmt"
)

type Addr struct {
	addr string
	name string
}

func NewAddr(addr, name string) *Addr {
	a := new(Addr)
	a.addr = addr
	a.name = name
	return a
}

func (a *Addr) Addr() string {
	return a.addr
}

func (a *Addr) Name() string {
	return a.name
}

func (a *Addr) String() string {
	if len(a.name) > 0 {
		return fmt.Sprintf("%s <%s>", encode(a.name), a.addr)
	} else {
		return fmt.Sprintf("<%s>", a.addr)
	}
}
