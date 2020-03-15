package campaignsrv

import "github.com/StephanHCB/go-campaign-service/internal/entity"

func validate(campaign *entity.Campaign) error {
	// some business validation

	// example: email address must not be @mailinator.com
	return nil
}