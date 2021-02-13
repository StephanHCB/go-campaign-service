package healthctl

import (
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/StephanHCB/go-campaign-service/web/util/timeout"
	"github.com/go-chi/chi"
	"net/http"
	"time"
)

func Create(server chi.Router) {
	server.Get("/health", Health)
	server.Get("/timeout", Timeout)
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// use this for easy mocking

var SleepTime = 10 * time.Second

func Timeout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	for i := 0; i < 10; i++ {
		aulogging.Logger.Ctx(ctx).Info().Printf("iteration %d - sleeping for 10 secs", i)
		time.Sleep(SleepTime)
		err := timeout.ErrIfTimeout(ctx)
		if err != nil {
			aulogging.Logger.Ctx(ctx).Warn().Print("timeout occurred")
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
