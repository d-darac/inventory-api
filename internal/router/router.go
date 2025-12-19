package router

import (
	"net/http"

	"github.com/d-darac/inventory-api/internal/groups"
	"github.com/d-darac/inventory-assets/api"
)

func LoadRoutes(mux *http.ServeMux, cfg *api.ApiConfig) {
	// TODO: Implement routes
	groupsHandler := groups.Handler{
		Db: cfg.Db,
	}

	mux.HandleFunc("POST /groups", groupsHandler.Create)
	mux.HandleFunc("DELETE /groups/{id}", groupsHandler.Delete)
	mux.HandleFunc("GET /groups", groupsHandler.List)
	mux.HandleFunc("GET /groups/{id}", groupsHandler.Retrieve)
	mux.HandleFunc("PUT /groups/{id}", groupsHandler.Update)
}
