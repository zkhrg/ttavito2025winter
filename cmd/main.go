package main

import (
	"database/sql"
	"log"

	"net/http"

	"ttavito/config"
	myHttp "ttavito/delivery/http"
	"ttavito/repository"
	"ttavito/usecase"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	repo := repository.NewEntityRepo(db)
	api := usecase.NewUsecase(repo)

	mux := http.NewServeMux()
	myHttp.SetupRoutes(api, mux)

	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
