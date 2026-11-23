package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mux := http.NewServeMux()
	
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
        logger.Info("health check", "remote", r.RemoteAddr)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}` + "\n"))
	})

	srv := &http.Server{
		Addr:              ":8081",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

}

