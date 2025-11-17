package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ethan-a-perry/song-loop/internal/database/dataaccess"
)

func main() {
	// Mongo
	db, err := dataaccess.NewMongoDataAccess()
	if err != nil {
		fmt.Println("could not get db")
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
