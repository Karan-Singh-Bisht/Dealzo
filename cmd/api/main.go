package main

import (
	"log"
	"net/http"
	"time"

	"github.com/karan-singh-bisht/Dealzo-api/internal/config"
)

func main() {

	cfg := config.MustLoad()

	// ServeMux routes incoming HTTP requests to their handlers.
	// Using an explicit mux makes the application's routing dependencies clear
	// and avoids relying on the package-level DefaultServeMux.
	mux := http.NewServeMux()

	// HTTP headers must be set before WriteHeader or Write.
	// If WriteHeader is not called explicitly, the first call to Write
	// implicitly sends a 200 OK response.
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// We need to manually set the server because there are several timeouts which are by default 0 meaning
	// the server waits for infinite time which should be only 30s or 60s see screenshot

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10, //default is 0
		WriteTimeout: time.Second * 30, //default is 0
		IdleTimeout:  time.Second * 60, //default is 0
	}

	log.Printf("server is listening on port %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server Not Connected %v", err)
	}
}
