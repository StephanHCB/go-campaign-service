package consumer

import (
	"context"
	"fmt"
	"github.com/StephanHCB/go-campaign-service/internal/repository/configuration"
	"github.com/StephanHCB/go-campaign-service/internal/repository/mailservice/mailserviceclient"
	"github.com/StephanHCB/go-campaign-service/web/util/media"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/middleware"
	"github.com/go-http-utils/headers"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"testing"
)

// contract test consumer side

const configKeySecuritySecret = "security.secret"
const configKeyDownstreamMailerserviceUrl = "downstream.mailerservice.url"

func TestConsumer(t *testing.T) {
	// Create Pact connecting to local Daemon
	pact := &dsl.Pact{
		Consumer: "CampaignService",
		Provider: "MailService",
		Host:     "localhost",
	}
	defer pact.Teardown()

	// types and values used in interaction
	type emailDto struct {
		ToAddress string `json:"to_address"`
		Subject   string `json:"subject"`
		Body      string `json:"body"`
	}
	tstToAddress := "demo@mailinator.com"
	tstSubject := "some subject"
	tstMailBody := "some mail body\nwith a second line\n"

	tstRequestId := "123456-request-id"

	// Pass in test case (consumer side)
	// This uses the repository on the consumer side to make the http call, should be as low level as possible
	var test = func() (err error) {
		// initialize test configuration so we have admin token key, mailer service url
		configuration.SetupForUnitTestDefaultsOnlyNoErrors()
		viper.Set(configKeySecuritySecret, "demosecret")
		viper.Set(configKeyDownstreamMailerserviceUrl, fmt.Sprintf("http://localhost:%d/", pact.Server.Port))

		// set up token and context
		tstSecret := configuration.SecuritySecret()
		tstValidationKeyFunc := func(token *jwt.Token) (interface{}, error) {
			return []byte(tstSecret), nil
		}
		token, err := jwt.Parse(tstAuthAdmin(), tstValidationKeyFunc)
		if err != nil {
			return err
		}

		ctx := context.TODO()
		// add specific authentication and Request ID to context
		ctx = context.WithValue(ctx, "user", token)
		ctx = context.WithValue(ctx, middleware.RequestIDKey, tstRequestId)

		client := mailserviceclient.Create()
		err = client.SendEmail(ctx, tstToAddress, tstSubject, tstMailBody)
		if err != nil {
			return err
		}
		// would test reply contents here if the client call returned a body
		return nil
	}

	// Set up our expected interactions.
	pact.
		AddInteraction().
		// contrived example, not really needed. This is the identifier of the state handler that will be called on the other side
		Given("an authorized user with the admin role exists").
		UponReceiving("A request to send an email").
		WithRequest(dsl.Request{
			Method: http.MethodPost,
			Headers: dsl.MapMatcher{
				headers.ContentType:   dsl.String(media.ContentTypeApplicationJson),
				headers.Authorization: dsl.String("Bearer " + tstAuthAdmin()),
			},
			Path: dsl.String("/api/rest/v1/sendmail"),
			Body: emailDto{
				ToAddress: tstToAddress,
				Subject:   tstSubject,
				Body:      tstMailBody,
			},
		}).
		WillRespondWith(dsl.Response{
			Status:  200,
			// Headers: dsl.MapMatcher{headers.ContentType: dsl.String(media.ContentTypeApplicationJson)},
			// Body:    dsl.String("OK"), (if we had a body for this request)
		})

	// Run the test, verify it did what we expected and capture the contract (writes a test log to logs/pact.log)
	if err := pact.Verify(test); err != nil {
		log.Fatalf("Error on Verify: %v", err)
	}

	// now write out the contract json (by default it goes to subdirectory pacts)
	if err := pact.WritePact(); err != nil {
		log.Fatalf("Error on pact write: %v", err)
	}
}
