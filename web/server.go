package web

import (
	"fmt"
	"github.com/StephanHCB/go-campaign-service/internal/repository/configuration"
	"github.com/StephanHCB/go-campaign-service/web/controller/swaggerctl"
	"github.com/StephanHCB/go-campaign-service/web/middleware/logfilter"
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

func Serve() {
	server := chi.NewRouter()

	server.Use(middleware.RequestID)

	server.Use(middleware.Logger)
	logfilter.Setup()

	swaggerctl.SetupSwaggerRoutes(server)

	address := configuration.ServerAddress()
	log.Info().Msg("Starting web server on " + address)
	err := http.ListenAndServe(address, server)
	if err != nil {
		// TODO log a warning on intentional shutdown, and an error otherwise
		failFunction(fmt.Errorf("Fatal error while starting web server: %s\n", err))
	}
}
