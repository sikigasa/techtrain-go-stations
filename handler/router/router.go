package router

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(db *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	healthzRouter(mux)
	todoRouter(mux, db)
	return mux
}

func healthzRouter(mux *http.ServeMux) *handler.HealthzHandler {
	mux.HandleFunc("/healthz", handler.NewHealthzHandler().ServeHTTP)
	return &handler.HealthzHandler{}
}

func todoRouter(mux *http.ServeMux, db *sql.DB) {
	todo := handler.NewTODOHandler(service.NewTODOService(db))
	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		var err error
		switch r.Method {
		case http.MethodPost:

			err = convertJson(todo.CreateTodo(w, r))
		// case http.MethodGet:
		case http.MethodPut:
			err = convertJson(todo.UpdateTodo(w, r))

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			convertJson(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			convertJson(w, http.StatusInternalServerError, err)
		}
	})
}

func convertJson(w http.ResponseWriter, status int, response interface{}) error {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}
	return nil
}
