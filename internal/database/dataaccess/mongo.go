package dataaccess

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongoDataAccess struct {
	Client *mongo.Client
	UserCollection *mongo.Collection
}

func NewMongoDataAccess() (*MongoDataAccess, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGO_CONNECTION_STRING")).SetServerAPIOptions(serverAPI)



	client, err := mongo.Connect(opts)
	if err != nil {
		fmt.Println("here")
		return nil, fmt.Errorf("Failed to create mongo client: %w", err)
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, fmt.Errorf("Could not send ping to confirm connection: %w", err)
	}

	fmt.Println("Pinged successfully.")

	return &MongoDataAccess{
		Client: client,
		UserCollection: client.Database(os.Getenv("DATABASE_NAME")).Collection(os.Getenv("USER_COLLECTION_NAME")),
	}, nil
}
