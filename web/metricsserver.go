package web

import (
	"github.com/StephanHCB/go-campaign-service/internal/repository/configuration"
	"github.com/rs/zerolog/log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// we use this with a go-routine, as it never returns
func metricsServerListenAndServe() {
	log.Info().Msg("starting metrics server on " + configuration.MetricsPort())
	metricsServeMux := http.NewServeMux()
	metricsServeMux.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(configuration.MetricsPort(), metricsServeMux)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start metrics http server: " + err.Error())
	}
}

func StartMetricsServer() {
	if configuration.EnableMetricsEndpoint() {
		go metricsServerListenAndServe()
	}
}
