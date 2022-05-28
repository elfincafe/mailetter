package mailetter

import (
	"io"
)

type Attachment struct {
	r		io.Reader
	name    string
}

func NewAttachment(r io.Reader, name string) (*Attachment, error) {
	a := new(Attachment)
	a.name = name
	a.r = r
	return a, nil
}
