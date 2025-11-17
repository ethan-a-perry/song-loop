package auth

import (
	"encoding/json"
	"net/http"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler {
		service: service,
	}
}

func (h *handler) GetUserFromRequest(w http.ResponseWriter, r *http.Request) {
	userId := GetContextValue(r, "user_id")

	user, err := h.service.GetUserFromRequest(userId)

	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)

		return
	}

	email := GetContextValue(r, "email")

	user, err = h.service.CreateUserFromRequest(userId, email)
	if err == nil {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "Authentication failed",
	})
}
