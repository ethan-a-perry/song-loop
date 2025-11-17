package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type SpotifyToken struct {
	AccessToken string
	TokenType string
	Scope string
	ExpiresAt time.Time
	RefreshToken string
}

type User struct {
	ID bson.ObjectID `bson:"_id"`
	Email string `bson:"email"`
	SpotifyToken SpotifyToken `bson:"spotify_token"`
}
