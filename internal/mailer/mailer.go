package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/wneessen/go-mail"
	"html/template"
	"time"
)

// templateFS holds the embed filesystem
// for the templates directory which contains
// the templates for sending emails.
//
//go:embed templates
var templateFS embed.FS

// Mailer is a struct which contains a mail.Client (used to
// connect to SMTP server and send mails) and sender information
// for your emails (the name and address you want the email to be from,
// such as "Alice Smith <alice@example.com>").
type Mailer struct {
	client *mail.Client
	sender string
}

func NewMailClient(host string, port int, username, password string) (*mail.Client, error) {
	return mail.NewClient(
		host,
		mail.WithPort(port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithTimeout(15*time.Second),
	)
}

func New(client *mail.Client, sender string) Mailer {
	return Mailer{client: client, sender: sender}
}

// Send is a method that takes the recipient email address as the first
// parameter, the name of the file containing the templates, and any dynamic
// data for the templates as an any parameter.
func (m *Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, fmt.Sprintf("templates/%s", templateFile))
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

	msg := mail.NewMsg()

	if err = msg.To(recipient); err != nil {
		return err
	}
	if err = msg.From(m.sender); err != nil {
		return err
	}

	msg.Subject(subject.String())
	msg.SetBodyString(mail.TypeTextPlain, plainBody.String())
	msg.AddAlternativeString(mail.TypeTextHTML, htmlBody.String())

	err = m.client.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}
