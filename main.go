package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/d-darac/inventory-api/internal/api"
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
		Host:       os.Getenv("HOST"),
		Port:       port,
		DbURL:      os.Getenv("DB_URL"),
		Platform:   os.Getenv("PLATFORM"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		MaxReqSize: 1024,
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

	router := http.NewServeMux()
	apiCfg.LoadRoutes(router)

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", router))

	stack := api.CreateStack(
		apiCfg.RecoveryMw,
		apiCfg.CheckReqBodyLengthMw,
		apiCfg.LoggerMw,
		apiCfg.ValidateJsonMw,
		apiCfg.CheckRouteAndMethodMw,
	)

	server := &http.Server{
		Handler: stack(v1.ServeHTTP),
		Addr:    fmt.Sprintf("%s:%d", apiCfg.Host, apiCfg.Port),
	}

	log.Printf("server listening on port: %d", apiCfg.Port)
	log.Fatal(server.ListenAndServe())
}
