package campaignsrv

import (
	"context"
	"errors"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/StephanHCB/go-campaign-service/internal/entity"
	"github.com/StephanHCB/go-campaign-service/internal/repository/database"
	"github.com/StephanHCB/go-campaign-service/internal/repository/database/dbrepo"
	"github.com/StephanHCB/go-campaign-service/internal/repository/mailservice"
)

type CampaignServiceImpl struct {
	DbRepository         dbrepo.Repository
	MailSenderRepository mailservice.MailSenderRepository
}

func Create(mailSender mailservice.MailSenderRepository) CampaignService {
	service := &CampaignServiceImpl{
		DbRepository: database.GetRepository(),
		MailSenderRepository: mailSender,
	}
	return service
}

func (s *CampaignServiceImpl) NewCampaign(ctx context.Context) *entity.Campaign {
	return &entity.Campaign{}
}

func (s *CampaignServiceImpl) CreateCampaign(ctx context.Context, campaign *entity.Campaign) (uint, error) {
	alreadyExists, err := s.isDuplicate(ctx, campaign.Subject, 0)
	if err != nil {
		return 0, err
	}
	if alreadyExists {
		aulogging.Logger.Ctx(ctx).Warn().Printf("received new campaign duplicate - subject: %v", campaign.Subject)
		return 0, errors.New("duplicate campaign data - same subject")
	}

	err = validate(campaign)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Warn().Printf("business validation for campaign failed - rejected: %v", err.Error())
		return 0, err
	}

	id, err := s.DbRepository.AddCampaign(ctx, campaign)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Printf("campaign write error: %v", err.Error())
	}
	return id, err
}

func (s *CampaignServiceImpl) UpdateCampaign(ctx context.Context, campaign *entity.Campaign) error {
	err := validate(campaign)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Warn().Printf("business validation for campaign update failed - changes rejected: %v", err.Error())
		return err
	}

	err = s.DbRepository.UpdateCampaign(ctx, campaign)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Printf("campaign write error: %v", err.Error())
	}
	return err
}

func (s *CampaignServiceImpl) GetCampaign(ctx context.Context, id uint) (*entity.Campaign, error) {
	campaign, err := s.DbRepository.GetCampaignById(ctx, id)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Printf("campaign read error: %v", err.Error())
	}
	return campaign, err
}

func (s *CampaignServiceImpl) ExecuteCampaign(ctx context.Context, campaign *entity.Campaign) (map[string]bool, error) {
	result := map[string]bool{}
	for _, recipient := range campaign.Recipients {
		aulogging.Logger.Ctx(ctx).Info().Printf("sending email subject '%s' to '%s'...", campaign.Subject, recipient.ToAddress)
		err := s.MailSenderRepository.SendEmail(ctx, recipient.ToAddress, campaign.Subject, campaign.Body)
		if err != nil {
			aulogging.Logger.Ctx(ctx).Error().WithErr(err).Printf("sending email subject '%s' to '%s' FAILED: %s", campaign.Subject, recipient.ToAddress, err.Error())
			result[recipient.ToAddress] = false
		} else {
			aulogging.Logger.Ctx(ctx).Info().Printf("sending email subject '%s' to '%s' successful", campaign.Subject, recipient.ToAddress)
			result[recipient.ToAddress] = true
		}
	}
	return result, nil
}

func (s *CampaignServiceImpl) isDuplicate(ctx context.Context, subject string, expectedCount uint) (bool, error) {
	count, err := s.DbRepository.CountCampaignsBySubject(ctx, subject)
	if err != nil {
		return false, err
	}
	return count != expectedCount, nil
}
