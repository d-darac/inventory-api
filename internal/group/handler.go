package group

import (
	"net/http"

	"github.com/d-darac/inventory-assets/database"
)

type Handler struct {
	Db *database.Queries
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) Retrieve(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {

}
