package healthctl

import (
	"github.com/go-chi/chi"
	"net/http"
)

func Create(server chi.Router) {
	server.Get("/health", Health)
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
