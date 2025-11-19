package models

import (
	"time"
)

type SpotifyToken struct {
	AccessToken string
	TokenType string
	Scope string
	ExpiresAt time.Time
	RefreshToken string
}
