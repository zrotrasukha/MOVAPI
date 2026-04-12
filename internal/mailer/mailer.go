package mailer

import (
	"bytes"
	"embed"
	"html/template"

	"gopkg.in/gomail.v2"
)

//go:embed templates
var templatesFS embed.FS

type Mailer struct {
	Dialer *gomail.Dialer
	Sender string
}

func New(host string, port int, username, password string, sender string) *Mailer {
	dialer := gomail.NewDialer(host, port, username, password)

	return &Mailer{
		Dialer: dialer,
		Sender: sender,
	}
}

func (m *Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templatesFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.Sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", subject.String())
	msg.AddAlternative("text/html", htmlBody.String())

	for range 3 {
		err = m.Dialer.DialAndSend(msg)
		if nil == err { // ain't no confusion widit
			return nil
		}
	}

	return err
}
