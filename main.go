package main

import (
	"log"
	"net/http"
)

const PORT = "8080"

func main() {
	servMux := http.ServeMux{}
	server := http.Server{
		Handler: &servMux,
		Addr:    ":" + PORT,
	}

	log.Printf("starting server on port: %s", PORT)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
