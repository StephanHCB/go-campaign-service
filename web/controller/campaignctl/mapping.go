package campaignctl

import (
	"github.com/StephanHCB/go-campaign-service/api/v1/campaign"
	"github.com/StephanHCB/go-campaign-service/internal/entity"
)

func mapDtoToCampaign(dto *campaign.CampaignDto, c *entity.Campaign) error {
	// do not map id - instead load by ID from db or you'll introduce errors
	c.Subject = dto.Subject
	c.Body	= dto.Body
	c.Recipients = []entity.Recipient{}
	for _, dtoRecipient := range dto.Recipients {
		r := entity.Recipient{}
		err := mapDtoToRecipient(&dtoRecipient, &r)
		if err != nil {
			return err
		}
		c.Recipients = append(c.Recipients, r)
	}
	return nil
}

func mapDtoToRecipient(dto *campaign.RecipientDto, r *entity.Recipient) error {
	// do not map id - instead load by ID from db
	r.ToAddress = dto.ToAddress
	return nil
}

func mapCampaignToDto(c *entity.Campaign, dto *campaign.CampaignDto)  {
	dto.Subject = c.Subject
	dto.Body = c.Body
	dto.Recipients = []campaign.RecipientDto{}
	for _, r := range c.Recipients {
		dtoRecipient := campaign.RecipientDto{}
		mapRecipientToDto(&r, &dtoRecipient)
		dto.Recipients = append(dto.Recipients, dtoRecipient)
	}
}

func mapRecipientToDto(r *entity.Recipient, dto *campaign.RecipientDto) {
	dto.ToAddress = r.ToAddress
}

func mapResultMapToExecutionResultDto(r map[string]bool, dto *campaign.ExecutionResultDto) {
	dto.Results = []campaign.SingleExecutionResultDto{}
	for k, v := range r {
		entry := campaign.SingleExecutionResultDto{
			ToAddress: k,
			Success:   v,
		}
		dto.Results = append(dto.Results, entry)
	}
}