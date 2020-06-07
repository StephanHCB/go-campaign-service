package mailservicemock

import (
	"context"
	"github.com/StephanHCB/go-campaign-service/internal/repository/mailservice"
	"github.com/rs/zerolog/log"
)

type MailSenderRepositoryMockImpl struct {
}

func Create() mailservice.MailSenderRepository {
	return &MailSenderRepositoryMockImpl{}
}

func (r *MailSenderRepositoryMockImpl) SendEmail(ctx context.Context, address string, subject string, body string) error {
	log.Ctx(ctx).Warn().Msgf("mock mailer SKIPPING call to mailer-service for address '%s', subject '%s', reporting successful send", address, subject)
	return nil
}
