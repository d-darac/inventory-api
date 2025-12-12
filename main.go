package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/d-darac/inventory-api/internal/api"
	"github.com/d-darac/inventory-api/internal/middleware"
	"github.com/d-darac/inventory-api/internal/router"
	"github.com/d-darac/inventory-assets/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("couldn't convert value of PORT env variable: %v", err)
	}

	apiCfg := api.ApiConfig{
		Host:     os.Getenv("HOST"),
		Port:     port,
		DbURL:    os.Getenv("DB_URL"),
		Platform: os.Getenv("PLATFORM"),
	}

	db, err := sql.Open("postgres", apiCfg.DbURL)
	if err != nil {
		log.Fatalf("couldn't open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("couldn't connect to database: %v", err)
	}

	apiCfg.Db = database.New(db)

	mux := http.NewServeMux()
	router.LoadRoutes(mux, &apiCfg)

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", mux))

	middleware := middleware.Middleware{
		MaxReqSize: 10240,
	}

	stack := middleware.CreateStack(
		middleware.RecoveryMw,
		middleware.CheckReqBodyLengthMw,
		middleware.LoggerMw,
		middleware.ValidateJsonMw,
		middleware.CheckRouteAndMethodMw,
	)

	server := &http.Server{
		Handler: stack(v1.ServeHTTP),
		Addr:    fmt.Sprintf("%s:%d", apiCfg.Host, apiCfg.Port),
	}

	log.Printf("server listening on port: %d", apiCfg.Port)
	log.Fatal(server.ListenAndServe())
}
