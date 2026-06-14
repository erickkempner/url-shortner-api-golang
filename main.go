package main

import (
	"log/slog"
	"net/http"
	"time"
	"urlshortener/api"
)

func main() {
	if err := run(); err != nil {
		slog.Error("failed to execute code", "error", err)
	}
}

func run() error {
	db := make(map[string]string)
	handler := api.NewHandler(db)

	s := http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: time.Minute,
		IdleTimeout:  10 * time.Second,
		Addr:         ":3000",
		Handler:      handler,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
