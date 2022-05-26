package mailetter

import "path/filepath"

type Attachment struct {
	content []byte
	name    string
	path    string
}

func NewAttachmentFrom(path string, name string) *Attachment {
	a := new(Attachment)
	a.path = path
	if len(name) > 0 {
		a.name = name
	} else {
		a.name = filepath.Base(path)
	}
	return a
}

func NewAttachment(name string, content []byte) *Attachment {
	a := new(Attachment)
	a.content = content
	a.name = name
	return a
}
