package auth

import (
	"fmt"
	"net/http"

	"github.com/ethan-a-perry/song-loop/internal/database/data"
)

type Service interface {
	Middleware(next http.Handler) http.Handler
	GetUserFromRequest(userId string) (UserResponse, error)
	CreateUserFromRequest(userId, email string) (UserResponse, error)
}

type svc struct {
	userData *data.UserData
}

type UserResponse struct {
	ID string `json:"id"`
	Email string `json:"email"`
}

func NewService(userData *data.UserData) Service {
	return &svc {
		userData: userData,
	}
}

func (s *svc) GetUserFromRequest(userId string) (UserResponse, error) {
	user, err := s.userData.GetUserById(userId)
	if err != nil {
		return UserResponse{}, fmt.Errorf("Could not find user by the id in database: %w", err)
	}

	return UserResponse{
		ID: userId,
		Email: user.Email,
	}, nil
}

func (s *svc) CreateUserFromRequest(userId, email string) (UserResponse, error) {
	if err := s.userData.CreateUser(userId, email); err != nil {
		return UserResponse{}, fmt.Errorf("Failed to create user: %w", err)
	}

	return UserResponse{
		ID: userId,
		Email: email,
	}, nil
}
