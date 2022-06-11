package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dhruv42/hraasaka/cache"
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

	err = db.CreateIndexOnHash(ctx, cfg, client)
	if err != nil {
		log.Fatalf("Error creating index: %v", err)
	}

	cache, err := cache.ConnectCache()
	if err != nil {
		log.Fatalf("Failed to create redis instance:- %v", err)
	}
	fmt.Println("********* REDIS INSTANCE CONNECTED SUCCESSFULLY ***********")

	val, err := cache.Exists("counter").Result()
	if err != nil {
		log.Fatalf("Failed to check if counter exists: %v", err)
	}

	if val == 0 {
		err = cache.Set("counter", 1000, 0).Err()
		if err != nil {
			log.Fatalf("Failed to set counter value: %v", err)
		}
	}
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
