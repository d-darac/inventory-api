package router

import (
	"net/http"

	"github.com/d-darac/inventory-api/internal/api"
	"github.com/d-darac/inventory-api/internal/group"
)

func LoadRoutes(mux *http.ServeMux, cfg *api.ApiConfig) {
	// TODO: Implement routes
	groupHandler := group.Handler{
		Db: cfg.Db,
	}

	mux.HandleFunc("POST /groups", groupHandler.Create)
	mux.HandleFunc("DELETE /groups/{id}", groupHandler.Delete)
	mux.HandleFunc("GET /groups", groupHandler.List)
	mux.HandleFunc("GET /groups/{id}", groupHandler.Retrieve)
	mux.HandleFunc("PUT /groups/{id}", groupHandler.Update)
}
