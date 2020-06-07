package requestidinresponse

import (
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func AddRequestIdHeaderToResponse(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestId := middleware.GetReqID(ctx)
		if requestId != "" {
			w.Header().Set(middleware.RequestIDHeader, requestId)
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}