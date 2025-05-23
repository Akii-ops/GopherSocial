package mailer

import "embed"

const (
	FromEmail           = "GopherSocial@sandbox.com"
	maxRetries          = 1
	UserWelcomeTemplate = "user_invitation.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any) error
}
