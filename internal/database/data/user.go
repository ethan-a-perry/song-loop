package data

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

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

func (ud *UserData) UpdateSpotifyToken(userId string, token *models.SpotifyToken) error {
	return nil
}
