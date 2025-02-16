package main

import (
	"log/slog"
	"time"

	"net/http"

	"ttavito/config"
	"ttavito/database"
	myHttp "ttavito/delivery/http"
	"ttavito/repository"
	"ttavito/usecase"
)

func main() {
	cfg := config.LoadConfig()

	pool, err := database.NewPostgresDB(cfg)
	if err != nil {
		slog.Error("Failed to create connection pool", "error", err)
		return
	}
	defer pool.Close()

	repo := repository.NewEntityRepo(pool)
	api := usecase.NewUsecase(repo)

	mux := http.NewServeMux()
	myHttp.SetupRoutes(api, mux)

	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	slog.Info("Server is running on port", "port", port)

	server := &http.Server{
		Addr:           ":" + port,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		slog.Error("Failed to start server: ", "error", err)
		return
	}
}
