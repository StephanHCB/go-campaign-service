package mailservice

import (
	"context"
)

type MailSenderRepository interface {
	SendEmail(ctx context.Context, address string, subject string, body string) error
}
