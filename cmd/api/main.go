package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	mux := http.NewServeMux()

	//Order of the header and write header matters
	//if we do not mention WriteHeader explicitly still Write method will call it on its own

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		num, err := w.Write([]byte(`{"status":"ok"}`))
		fmt.Println(num)
		if err != nil {
			log.Fatalf("Dealzo API is not running...")
		}
	})

	// We need to manually set the server because there are several timeouts which are by default 0 meaning
	// the server waits for infinite time which should be only 30s or 60s see screenshot

	srv := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server Not Connected %v", err)
	}
}
