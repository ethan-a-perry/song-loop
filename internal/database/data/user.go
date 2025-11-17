package data

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	// "go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/ethan-a-perry/song-loop/internal/database/models"
)

type UserData struct {
	users *mongo.Collection
}

func NewUserData(users *mongo.Collection) *UserData {
	return &UserData{
		users: users,
	}
}

func (ud *UserData) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	if err := ud.users.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, fmt.Errorf("Failed to find user document: %w", err)
	}

	return &user, nil
}

// func (ud *UserData) GetSpotifyToken() (*Token, error) {
// 	return "", nil
// }

// func (ud *UserData) UpdateSpotifyToken(token *Token) error {
// 	filter := bson.M{"_id": ""}

// 	update := bson.M{
//         "$set": bson.M{
//             "spotify_token": &token,
//         },
//     }

// 	opts := options.Replace().SetUpsert(true)

// 	_, err := ud.users.ReplaceOne(context.TODO(), filter, update, opts)

// 	if err != nil {
// 		return fmt.Errorf("failed to save spotify token: %w", err)
// 	}

// 	return nil
// }

func (ud *UserData) CreateUser(id, email string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	user := models.User{
		ID: objectID,
		Email: email,
	}

	_, err = ud.users.InsertOne(context.TODO(), user)
	if err != nil {
		return fmt.Errorf("Failed to insert user document: %w", err)
	}

	return nil
}

func (ud *UserData) GetUserById(id string) (*models.User, error) {
	var user models.User

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	if err := ud.users.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user); err != nil {
		return nil, fmt.Errorf("Failed to find user document: %w", err)
	}

	return &user, nil
}
