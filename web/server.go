package web

import (
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/StephanHCB/go-autumn-logging-zerolog/loggermiddleware"
	"github.com/StephanHCB/go-campaign-service/internal/repository/configuration"
	"github.com/StephanHCB/go-campaign-service/internal/service/campaignsrv"
	"github.com/StephanHCB/go-campaign-service/web/controller/campaignctl"
	"github.com/StephanHCB/go-campaign-service/web/controller/healthctl"
	"github.com/StephanHCB/go-campaign-service/web/controller/swaggerctl"
	"github.com/StephanHCB/go-campaign-service/web/middleware/authentication"
	"github.com/StephanHCB/go-campaign-service/web/middleware/requestidinresponse"
	"github.com/StephanHCB/go-campaign-service/web/middleware/requestlogging"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"time"
)

// use this for easy mocking

var failFunction = fail

var RequestTimeout = 2 * time.Minute

func fail(err error) {
	// this does os.exit 1
	aulogging.Logger.NoCtx().Fatal().WithErr(err).Print(err.Error())
}

func Create(campaignService campaignsrv.CampaignService) chi.Router {
	aulogging.Logger.NoCtx().Info().Print("Creating router and setting up filter chain")
	server := chi.NewRouter()

	middleware.RequestIDHeader = "X-B3-TraceId"
	server.Use(middleware.RequestID)

	server.Use(middleware.Logger)
	requestlogging.Setup()

	server.Use(middleware.Timeout(RequestTimeout))

	server.Use(loggermiddleware.AddZerologLoggerToContext)

	server.Use(requestidinresponse.AddRequestIdHeaderToResponse)

	server.Use(authentication.AddJWTTokenInfoToContext(configuration.SecuritySecret()))

	aulogging.Logger.NoCtx().Info().Print("Setting up routes")

	_ = campaignctl.Create(server, campaignService)

	healthctl.Create(server)

	swaggerctl.SetupSwaggerRoutes(server)

	return server
}

func Serve(server chi.Router) {
	address := configuration.ServerAddress()
	aulogging.Logger.NoCtx().Info().Print("Starting web server on " + address)

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		// TLSConfig:    tlsConfig,
		Addr:    address,
		Handler: server,
	}
	err := srv.ListenAndServe()
	if err != nil {
		// TODO log a warning on intentional shutdown, and an error otherwise
		failFunction(fmt.Errorf("Fatal error while starting web server: %s\n", err))
	}
}
