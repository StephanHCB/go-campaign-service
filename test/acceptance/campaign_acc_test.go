package acceptance

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

// ------------------------------------------
// acceptance tests for the attendee resource
// ------------------------------------------

// --- create new attendee ---

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
				So(response.location, shouldMatchRegex, "^\\/api\\/rest\\/v1\\/campaigns\\/[1-9][0-9]*$")

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
