package web

import (
	"fmt"
	"github.com/StephanHCB/go-campaign-service/web/controller/swaggerctl"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

// use this for easy mocking

var failFunction = fail

func fail(err error) {
	panic(err)
}

func Serve() {
	server := chi.NewRouter()
	server.Use(middleware.Logger)

	swaggerctl.SetupSwaggerRoutes(server)

	err := http.ListenAndServe(":8081", server)
	if err != nil {
		// TODO log a warning on intentional shutdown, and an error otherwise
		failFunction(fmt.Errorf("Fatal error while starting web server: %s\n", err))
	}
}
