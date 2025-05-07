package mailer

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
	"time"

	"github.com/go-mail/mail"
)

type EtherealMailer struct {
	fromEmail string

	client *mail.Dialer
}

func NewEtherealMailer(fromEmail, host, username, password string, port int) *EtherealMailer {
	client := mail.NewDialer(host, port, username, password)

	// 强制TLS
	client.StartTLSPolicy = mail.MandatoryStartTLS
	s, err := client.Dial()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer s.Close()

	return &EtherealMailer{
		fromEmail: fromEmail,
		client:    client,
	}
}

func (m *EtherealMailer) Send(templateFile, username, email string, data any) error {

	// template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+UserWelcomeTemplate)
	if err != nil {
		return err
	}
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}
	msg := mail.NewMessage()
	msg.SetHeader("From", m.fromEmail)
	msg.SetHeader("To", email)
	msg.SetAddressHeader("Cc", "GopherSocial@test.com", "Aki")
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/html", body.String())

	var retryErr error
	for i := 0; i < maxRetries; i++ {
		retryErr = m.client.DialAndSend(msg)
		if retryErr != nil {
			// backoff
			// log.Printf("Failed to send email to %v, attempt %d of %d\n", email, i+1, maxRetries)
			log.Printf("test------------------- : %v\n", err)
			time.Sleep(time.Second * time.Duration(i+1))

			continue
		}

		log.Printf("Email sent with status code 200")

		return nil

	}
	return fmt.Errorf("failed to send email after %d attempts,error: %v", maxRetries, retryErr)

}
