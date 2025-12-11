package api

import "github.com/d-darac/inventory-assets/database"

type ApiConfig struct {
	Host     string
	Port     int
	DbURL    string
	Db       *database.Queries
	Platform string
}
