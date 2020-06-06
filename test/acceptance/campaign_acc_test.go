package acceptance

import (
	"github.com/StephanHCB/go-campaign-service/api/v1/campaign"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

// ------------------------------------------
// acceptance tests for the campaign resource
// ------------------------------------------

// --- create new campaign ---

func TestCreateNewCampaignAsAdmin(t *testing.T) {
	tstSetup(tstValidConfigurationPath)
	defer tstShutdown()

	Convey("Given an authenticated user with the admin role", t, func() {
		token := tstAuthAdmin()

		Convey("When they try to create a campaign", func() {
			campaignSent := tstBuildValidCampaign("campaign-create-admin")
			response, err := tstPerformPut("/api/rest/v1/campaigns", tstRenderJson(campaignSent), token)

			Convey("Then the campaign is successfully created", func() {
				So(err, ShouldEqual, nil)
				So(response.status, ShouldEqual, http.StatusCreated)
				So(response.location, ShouldMatchRegex, "^\\/api\\/rest\\/v1\\/campaigns\\/[1-9][0-9]*$")

				Convey("And the same campaign can be read again", func() {
					campaignReadAgain, err := tstReadCampaign(response.location)
					So(err, ShouldEqual, nil)
					So(campaignReadAgain, ShouldDeepEqual, campaignSent)
				})
			})
		})
	})
}

// security acceptance tests examples

func TestCreateNewCampaignUnauthenticated_Deny(t *testing.T) {
	tstSetup(tstValidConfigurationPath)
	defer tstShutdown()

	Convey("Given an unauthenticated user", t, func() {
		token := ""

		Convey("When they try to create a campaign", func() {
			campaignSent := tstBuildValidCampaign("campaign-create-unauth")
			response, err := tstPerformPut("/api/rest/v1/campaigns", tstRenderJson(campaignSent), token)

			Convey("Then the request is denied as unauthenticated", func() {
				So(err, ShouldEqual, nil)
				So(response.status, ShouldEqual, http.StatusUnauthorized)
				So(response.location, ShouldEqual, "")
			})
		})
	})
}

func TestCreateNewCampaignUnauthorized_Deny(t *testing.T) {
	tstSetup(tstValidConfigurationPath)
	defer tstShutdown()

	Convey("Given an authenticated user that does not have the admin role", t, func() {
		token := tstAuthUser()

		Convey("When they try to create a campaign", func() {
			campaignSent := tstBuildValidCampaign("campaign-create-unauth")
			response, err := tstPerformPut("/api/rest/v1/campaigns", tstRenderJson(campaignSent), token)

			Convey("Then the request is denied as unauthorized", func() {
				So(err, ShouldEqual, nil)
				So(response.status, ShouldEqual, http.StatusForbidden)
				So(response.location, ShouldEqual, "")
			})
		})
	})
}

// --- execute campaign ---

func TestExecuteCampaignAsAdmin(t *testing.T) {
	tstSetup(tstValidConfigurationPath)
	defer tstShutdown()

	Convey("Given an authenticated user with the admin role", t, func() {
		token := tstAuthAdmin()

		Convey("And an existing campaign", func() {
			campaignSent := tstBuildValidCampaign("campaign-execute-admin")
			createResponse, err := tstPerformPut("/api/rest/v1/campaigns", tstRenderJson(campaignSent), token)
			So(err, ShouldEqual, nil)
			So(createResponse.status, ShouldEqual, http.StatusCreated)
			So(createResponse.location, ShouldMatchRegex, "^\\/api\\/rest\\/v1\\/campaigns\\/[1-9][0-9]*$")

			Convey("When they try to execute the campaign", func() {
				response, err := tstPerformPost(createResponse.location + "/execute", "", token)

				Convey("Then the emails are sent out", func() {
					expectedResponse := campaign.ExecutionResultDto{Results: []campaign.SingleExecutionResultDto{
						{ToAddress: tstEmailAddress, Success: true},
					}}
					actualResponseBody := campaign.ExecutionResultDto{}
					parseErr := tstParseJson(response.body, &actualResponseBody)

					// TODO validate contents of email

					So(err, ShouldEqual, nil)
					So(response.status, ShouldEqual, http.StatusOK)
					So(response.location, ShouldEqual, createResponse.location)
					So(parseErr, ShouldEqual, nil)
					So(actualResponseBody, ShouldDeepEqual, expectedResponse)
				})
			})
		})
	})
}
