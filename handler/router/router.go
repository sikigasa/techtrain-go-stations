package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
)

// struct HealthzResponse {

// }

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	healthzRouter(mux)
	return mux
}

func healthzRouter(mux *http.ServeMux) *handler.HealthzHandler {
	mux.HandleFunc("/healthz", handler.NewHealthzHandler().ServeHTTP)
	return &handler.HealthzHandler{}
}
