package mailservicemock

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/StephanHCB/go-campaign-service/internal/repository/mailservice"
)

type MailSenderRepositoryMockImpl struct {
}

func Create() mailservice.MailSenderRepository {
	return &MailSenderRepositoryMockImpl{}
}

func (r *MailSenderRepositoryMockImpl) SendEmail(ctx context.Context, address string, subject string, body string) error {
	aulogging.Logger.Ctx(ctx).Warn().Printf("mock mailer SKIPPING call to mailer-service for address '%s', subject '%s', reporting successful send", address, subject)
	return nil
}
