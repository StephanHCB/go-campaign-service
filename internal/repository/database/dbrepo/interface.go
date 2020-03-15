package dbrepo

import (
	"context"
	"github.com/StephanHCB/go-campaign-service/internal/entity"
)

type Repository interface {
	Open()
	Close()
	Migrate()

	AddCampaign(ctx context.Context, a *entity.Campaign) (uint, error)
	UpdateCampaign(ctx context.Context, a *entity.Campaign) error
	GetCampaignById(ctx context.Context, id uint) (*entity.Campaign, error)
}
