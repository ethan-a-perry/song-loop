package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ethan-a-perry/song-loop/internal/database/dataaccess"
	"github.com/joho/godotenv"
)

func main() {
	// ENV file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	// Mongo
	db, err := dataaccess.NewMongoDataAccess()
	if err != nil {
		fmt.Println("could not get db")
		fmt.Println(err)
	}

	defer func() {
		// Disconnect mongo client
		if err := db.Client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Configure API
	cfg := config {
		addr: ":8080",
	}

	api := api {
		config: cfg,
		db: db,
	}

	router := api.mount()

	if err := api.run(router); err != nil {
		fmt.Printf("Server has failed to start: %s", err)
		os.Exit(1)
	}
}
