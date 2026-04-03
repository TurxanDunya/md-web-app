package main

import (
	"log"
	"md_api/internal/config"
	"md_api/internal/database"
	"md_api/internal/routes"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	pool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	mux := http.NewServeMux()
	routes.Setup(mux, pool, cfg)

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
