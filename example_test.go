package mailetter_test

import (
	"fmt"
	"mailetter"
	"os"
)

func Example() {
	// Creating Mail
	from := mailetter.NewAddr("from@example.com", "Example Sender")
	m := mailetter.NewMail(from)
	m.Header("X-Mailer", "MaiLetter Mail Client")
	to := mailetter.NewAddr("to@example.com", "Recipient User")
	m.To(to)
	cc := mailetter.NewAddr("cc@example.com", "Carbon Copy Recipient ")
	m.Cc(cc)
	bcc := mailetter.NewAddr("bcc@example.com", "Blind Carbon Copy Recipient ")
	m.Bcc(bcc)
	m.Subject("Example Subject")
	body := `
Dear {{.Name}}

This is an example mail.
`
	m.Body(body)
	m.Set("Name", "Recipient")
	// Connecting and Sending mail
	ml, _ := mailetter.New("smtps://mail.example.com:465")
	err := ml.Send(m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}
}
