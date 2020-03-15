package ctxlogger

import (
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
	"net/http"
)

func AddZerologLoggerToContext(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestId := middleware.GetReqID(ctx)
		method := r.Method
		path := r.URL.Path

		sublogger := log.Logger.With().
			Str("request-id", requestId).
			Str("method", method).
			Str("path", path).
			Logger()
		newCtx := sublogger.WithContext(ctx)

		next.ServeHTTP(w, r.WithContext(newCtx))
	}
	return http.HandlerFunc(fn)
}