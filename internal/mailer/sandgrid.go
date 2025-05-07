package mailer

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"text/template"
	"time"

	sendcloud "github.com/sendcloud2013/sendcloud-sdk-go/email"
)

// 不导出
type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendcloud.SendCloud
}

func NewSendGrid(apiKey, fromEmail string) *SendGridMailer {
	client, err := sendcloud.NewSendCloud(fromEmail, apiKey)
	if err != nil {
		log.Fatal(err)
	}
	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandBox bool) error {
	_ = isSandBox
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
	args := &sendcloud.CommonMail{
		Receiver: sendcloud.MailReceiver{
			To: email,
		},
		Body: sendcloud.MailBody{
			From:     m.fromEmail,
			Subject:  subject.String(),
			FromName: FromEmail,
		},
		Content: sendcloud.TextContent{
			Html: body.String(),
		},
	}
	// from := mail.NewEmail(FromEmail, m.fromEmail)

	// to := mail.NewEmail(username, email)

	// message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	// set sandbox mode
	// message.SetMailSettings(&mail.MailSettings{
	// 	SandboxMode: &mail.Setting{
	// 		Enable: &isSandBox,
	// 	},
	// })

	ctx := context.Background()

	for i := 0; i < maxRetries; i++ {
		response, err := m.client.SendCommonEmail(ctx, args)
		if err != nil {
			log.Printf("Failed to send email to %v, attempt %d of %d\n", email, i+1, maxRetries)
			log.Printf("Error 22 : %v\n", err)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		log.Printf("Email sent with status code %v", response.StatusCode)
	}
	return fmt.Errorf("failed to send email after %d attempts", maxRetries)

}
