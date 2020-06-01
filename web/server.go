package web

import (
	"fmt"
	"github.com/StephanHCB/go-campaign-service/internal/repository/configuration"
	"github.com/StephanHCB/go-campaign-service/internal/service/campaignsrv"
	"github.com/StephanHCB/go-campaign-service/web/controller/campaignctl"
	"github.com/StephanHCB/go-campaign-service/web/controller/healthctl"
	"github.com/StephanHCB/go-campaign-service/web/controller/swaggerctl"
	"github.com/StephanHCB/go-campaign-service/web/middleware/authentication"
	"github.com/StephanHCB/go-campaign-service/web/middleware/ctxlogger"
	"github.com/StephanHCB/go-campaign-service/web/middleware/requestlogging"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
	"net/http"
)

// use this for easy mocking

var failFunction = fail

func fail(err error) {
	// this does os.exit 1
	log.Fatal().Err(err).Msg(err.Error())
}

func Create() chi.Router {
	log.Info().Msg("Creating router and setting up filter chain")
	server := chi.NewRouter()

	server.Use(middleware.RequestID)

	server.Use(middleware.Logger)
	requestlogging.Setup()

	server.Use(ctxlogger.AddZerologLoggerToContext)

	server.Use(authentication.AddJWTTokenInfoToContext(configuration.SecuritySecret()))

	log.Info().Msg("Setting up routes")

	_ = campaignctl.Create(server, campaignsrv.Create())

	healthctl.Create(server)

	swaggerctl.SetupSwaggerRoutes(server)

	return server
}

func Serve(server chi.Router) {
	address := configuration.ServerAddress()
	log.Info().Msg("Starting web server on " + address)
	err := http.ListenAndServe(address, server)
	if err != nil {
		// TODO log a warning on intentional shutdown, and an error otherwise
		failFunction(fmt.Errorf("Fatal error while starting web server: %s\n", err))
	}
}
