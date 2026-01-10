package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
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
		log.Fatalf("[main] Couldn't convert value of PORT env variable to int: %v", err)
	}

	apiCfg := api.ApiConfig{
		DbURL:    env.DB_URL,
		Host:     env.HOST,
		Platform: env.PLATFORM,
		Port:     port,
	}

	db, err := sql.Open("postgres", apiCfg.DbURL)
	if err != nil {
		log.Fatalf("[main] Couldn't open database connection: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("[main] Couldn't ping to database: %v", err)
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
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}

	go func() {
		log.Print("[main] Server starting...")
		var err error
		if strings.ToLower(apiCfg.Platform) == "dev" {
			err = server.ListenAndServe()
		} else {
			err = server.ListenAndServeTLS(env.TLS_CERT_PATH, env.TLS_KEY_PATH)
		}
		if err != http.ErrServerClosed {
			log.Printf("[main] Failed to start server: %v.", err)
		} else {
			log.Printf("[main] %v.", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	fmt.Println()
	log.Print("[main] Server stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println()
		log.Printf("[main] Server forced to shutdown: %v.", err)
		return
	}

	log.Print("[main] Server stopped.")
}
