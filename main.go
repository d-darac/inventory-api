package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/d-darac/inventory-api/env"
	"github.com/d-darac/inventory-api/middleware"
	"github.com/d-darac/inventory-api/router"
	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	_ "github.com/lib/pq"
)

func main() {
	env := env.GetEnv()

	port, err := strconv.Atoi(env.PORT)
	if err != nil {
		log.Fatalf("couldn't convert value of PORT env variable: %v", err)
	}

	apiCfg := api.ApiConfig{
		DbURL:    env.DB_URL,
		Host:     env.HOST,
		Platform: env.PLATFORM,
		Port:     port,
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
		Db:         apiCfg.Db,
		Auth: struct {
			MasterKey string
			Iv        string
		}{
			MasterKey: env.MASTER_KEY,
			Iv:        env.IV,
		},
	}

	stack := middleware.CreateStack(
		middleware.RecoveryMw,
		middleware.CheckReqBodyLengthMw,
		middleware.LoggerMw,
		middleware.CheckRouteAndMethodMw,
		middleware.ApiKeyAuthMw,
	)

	server := &http.Server{
		Handler:           stack(v1.ServeHTTP),
		Addr:              fmt.Sprintf("%s:%d", apiCfg.Host, apiCfg.Port),
		ReadHeaderTimeout: time.Second * 15,
	}

	log.Printf("server listening on port: %d", apiCfg.Port)
	log.Fatal(server.ListenAndServe())
}
