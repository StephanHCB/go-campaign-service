package campaignctl

import (
	"context"
	"github.com/StephanHCB/go-campaign-service/api/v1/campaign"
)

func validate(ctx context.Context, dto *campaign.CampaignDto) []string {
	// some syntactical validation, but not business rules
	return []string{}
}
