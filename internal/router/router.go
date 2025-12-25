package router

import (
	"net/http"

	"github.com/d-darac/inventory-api/internal/groups"
	"github.com/d-darac/inventory-api/internal/inventories"
	itemidentifiers "github.com/d-darac/inventory-api/internal/item_identifiers"
	"github.com/d-darac/inventory-api/internal/items"
	"github.com/d-darac/inventory-api/internal/prices"
	"github.com/d-darac/inventory-assets/api"
)

func LoadRoutes(mux *http.ServeMux, cfg *api.ApiConfig) {
	// TODO: Implement routes
	groupsHandler := groups.NewHandler(cfg.Db)
	mux.HandleFunc("POST /groups", groupsHandler.Create)
	mux.HandleFunc("DELETE /groups/{id}", groupsHandler.Delete)
	mux.HandleFunc("GET /groups", groupsHandler.List)
	mux.HandleFunc("GET /groups/{id}", groupsHandler.Retrieve)
	mux.HandleFunc("PUT /groups/{id}", groupsHandler.Update)

	itemsHandler := items.NewHandler(cfg.Db)
	mux.HandleFunc("POST /items", itemsHandler.Create)
	mux.HandleFunc("DELETE /items/{id}", itemsHandler.Delete)
	mux.HandleFunc("GET /items", itemsHandler.List)
	mux.HandleFunc("GET /items/{id}", itemsHandler.Retrieve)
	mux.HandleFunc("PUT /items/{id}", itemsHandler.Update)

	inventoriesHandler := inventories.NewHandler(cfg.Db)
	mux.HandleFunc("POST /inventories", inventoriesHandler.Create)
	mux.HandleFunc("DELETE /inventories/{id}", inventoriesHandler.Delete)
	mux.HandleFunc("GET /inventories", inventoriesHandler.List)
	mux.HandleFunc("GET /inventories/{id}", inventoriesHandler.Retrieve)
	mux.HandleFunc("PUT /inventories/{id}", inventoriesHandler.Update)

	itemIdentifiersHandler := itemidentifiers.NewHandler(cfg.Db)
	mux.HandleFunc("POST /item_identifiers", itemIdentifiersHandler.Create)
	mux.HandleFunc("DELETE /item_identifiers/{id}", itemIdentifiersHandler.Delete)
	mux.HandleFunc("GET /item_identifiers", itemIdentifiersHandler.List)
	mux.HandleFunc("GET /item_identifiers/{id}", itemIdentifiersHandler.Retrieve)
	mux.HandleFunc("PUT /item_identifiers/{id}", itemIdentifiersHandler.Update)

	prices := prices.NewHandler(cfg.Db)
	mux.HandleFunc("POST /prices", prices.Create)
	mux.HandleFunc("DELETE /prices/{id}", prices.Delete)
	mux.HandleFunc("GET /prices", prices.List)
	mux.HandleFunc("GET /prices/{id}", prices.Retrieve)
	mux.HandleFunc("PUT /prices/{id}", prices.Update)
}
