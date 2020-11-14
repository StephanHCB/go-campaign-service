package web

import (
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/StephanHCB/go-campaign-service/internal/repository/configuration"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// we use this with a go-routine, as it never returns
func metricsServerListenAndServe() {
	aulogging.Logger.NoCtx().Info().Print("starting metrics server on " + configuration.MetricsPort())
	metricsServeMux := http.NewServeMux()
	metricsServeMux.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(configuration.MetricsPort(), metricsServeMux)
	if err != nil {
		aulogging.Logger.NoCtx().Fatal().WithErr(err).Print("failed to start metrics http server: " + err.Error())
	}
}

func StartMetricsServer() {
	if configuration.EnableMetricsEndpoint() {
		go metricsServerListenAndServe()
	}
}
