package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Siddharth9101/olx-api/internal/config"
)

func main() {
	// loading config
	cfg := config.MustLoad()
	// initializing new router/mux
	mux := http.NewServeMux();

	// defining health endpoint
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	// initializing new server
	srv := http.Server{
		Addr: ":"+ cfg.Port,
		Handler: mux,
		ReadTimeout: time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout: time.Second * 60,
	}
	log.Println("server is listning on "+srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}