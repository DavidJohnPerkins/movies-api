package main

import (
	"context"
	"dperkins/movies-api/api"
	"dperkins/movies-api/config"
	"dperkins/movies-api/store"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	store := store.NewSqlServerMoviesStore(cfg.DatabaseURL)
	server := api.NewServer(cfg.HTTPServer, store)
	server.Start(ctx)
}
