package campaign

import "net/http"

// --- models ---

// Model for CampaignDto.
//
// swagger:model campaignDto
type CampaignDto struct {
	// The email subject
	Subject string `json:"subject"`
	// The email body
	Body string `json:"body"`
	// The list of recipients
	Recipients []RecipientDto `json:"recipients"`
}

// Model for RecipientDto.
//
// swagger:model recipientDto
type RecipientDto struct {
	// The email address to send to
	ToAddress string `json:"to_address"`
}

// Model for ExecutionResultDto.
//
// swagger:model executionResultDto
type ExecutionResultDto struct {
	Results []SingleExecutionResultDto `json:"results"`
}

// Model for SingleExecutionResultDto.
//
// swagger:model singleExecutionResultDto
type SingleExecutionResultDto struct {
	// The email address to send to
	ToAddress string `json:"to_address"`

	// Whether sending was successful
	Success bool `json:"success"`
}

// --- parameters and responses --- needed to use models

// Parameters for creating a Campaign
//
// swagger:parameters createCampaignParams
type CreateCampaignParams struct {
	// in:body
	Body CampaignDto
}

// Parameters for updating a Campaign
//
// swagger:parameters updateCampaignParams
type UpdateCampaignParams struct {
	// The id of the campaign
	//
	// in:path
	Id string

	// The changed data of the campaign to be set
	//
	// in:body
	Body CampaignDto
}

// Parameters for fetching or addressing a Campaign
//
// swagger:parameters addressCampaignParams
type AddressCampaignParams struct {
	// The id of the campaign
	//
	// in:path
	Id string
}

// The campaign location response with just a Location header
//
// swagger:response campaignLocationResponse
type CampaignLocationResponse struct {
	// Location header
	Location string `json:"Location"`
}

// The campaign data response including a Location header
//
// swagger:response campaignDataResponse
type CampaignDataResponse struct {
	// Location header
	Location string `json:"Location"`

	// The data of the campaign
	//
	// in:body
	Body CampaignDto
}

// The campaign execution response including a Location header and the list of addresses
//
// swagger:response campaignExecutionResponse
type CampaignExecutionResponse struct {
	// Location header
	Location string `json:"Location"`

	// The data of the campaign
	//
	// in:body
	Body ExecutionResultDto
}

// --- routes ---

type CampaignApi interface {
	// swagger:route PUT /api/rest/v1/campaigns campaign-tag createCampaignParams
	// This will create a new campaign and return its location
	//
	// responses:
	//   201: campaignLocationResponse
	//   400: errorResponse
	//   401: errorResponse
	//   403: errorResponse
	CreateCampaign(w http.ResponseWriter, r *http.Request)

	// swagger:route POST /api/rest/v1/campaigns/{Id} campaign-tag updateCampaignParams
	// This will update an existing campaign
	//
	// responses:
	//   200: campaignLocationResponse
	//   400: errorResponse
	//   401: errorResponse
	//   403: errorResponse
	//   404: errorResponse
	UpdateCampaign(w http.ResponseWriter, r *http.Request)

	// swagger:route GET /api/rest/v1/campaigns/{Id} campaign-tag addressCampaignParams
	// This will return an existing campaign
	//
	// responses:
	//   200: campaignDataResponse
	//   401: errorResponse
	//   403: errorResponse
	//   404: errorResponse
	GetCampaign(w http.ResponseWriter, r *http.Request)

	// swagger:route POST /api/rest/v1/campaigns/{Id}/execute campaign-tag addressCampaignParams
	// This will execute a campaign, that is send emails using the mailer service.
	//
	// responses:
	//   200: campaignExecutionResponse
	//   400: errorResponse
	//   401: errorResponse
	//   403: errorResponse
	//   404: errorResponse
	//   502: errorResponse
	ExecuteCampaign(w http.ResponseWriter, r *http.Request)
}
