package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dhruv42/hraasaka/config"
	"github.com/dhruv42/hraasaka/db"
	"github.com/dhruv42/hraasaka/routes"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Error loading config: [%v]", err)
	}

	client, err := db.GetDBConnection(ctx, cfg)
	if err != nil {
		log.Fatalf("Db connectiom failed:- %v", err)
	}

	fmt.Println("***** Database connected *****")
	defer client.Disconnect(ctx)

	// Start the server
	serve(cfg)
}

func serve(cfg *config.Config) {
	fmt.Printf("========== Server listening on port %d ==========\n", cfg.Port)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: routes.New(),
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
