package campaignsrv

import (
	"context"
	"github.com/StephanHCB/go-campaign-service/internal/entity"
)

type CampaignService interface {
	NewCampaign(ctx context.Context) *entity.Campaign

	CreateCampaign(ctx context.Context, campaign *entity.Campaign) (uint, error)
	UpdateCampaign(ctx context.Context, campaign *entity.Campaign) error
	GetCampaign(ctx context.Context, id uint) (*entity.Campaign, error)

	ExecuteCampaign(ctx context.Context, campaign *entity.Campaign) (map[string]bool, error)
}
