package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Siddharth9101/olx-api/internal/config"
	"github.com/Siddharth9101/olx-api/internal/db"
	"github.com/Siddharth9101/olx-api/internal/handlers"
)

func main() {
	// loading config
	cfg := config.MustLoad()

	// configuring db
	db, err := db.Connect(cfg.DatabaseUrl)
	if err != nil{
		log.Fatalf("main.db.connect: %v", err)
	}

	// initializing new router/mux
	mux := http.NewServeMux();

	lh := handlers.NewListingHandler(db)

	// api endpoints
	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("GET /listings", lh.List)
	mux.HandleFunc("DELETE /listings/{id}", lh.Delete)

	// initializing new server
	srv := http.Server{
		Addr: ":"+ cfg.Port,
		Handler: mux,
		ReadTimeout: time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout: time.Second * 60,
	}
	log.Println("database connected")
	log.Println("server is listning on "+srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}