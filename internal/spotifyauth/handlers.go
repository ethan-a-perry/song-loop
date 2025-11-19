package spotifyauth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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

func (h *handler) Connect(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// base64 encode the userID
	state := base64.URLEncoding.EncodeToString([]byte(userID))

	authorizationUrl, err := h.service.GetAuthorizationUrl(state)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"authUrl": authorizationUrl,
	})
}

func (h *handler) Callback(w http.ResponseWriter, r *http.Request) {
	errMsg := r.URL.Query().Get("error")
	if errMsg != "" {
		http.Error(w, "Authorization failed during callback: " + errMsg, http.StatusUnauthorized)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code provided", http.StatusBadRequest)
        return
	}

	stateEncoded := r.URL.Query().Get("state")
	if stateEncoded == "" {
		http.Error(w, "No state provided", http.StatusBadRequest)
		return
	}

	userIDBytes, err := base64.URLEncoding.DecodeString(stateEncoded)
	if err != nil {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}
	userID := string(userIDBytes)

	if err := h.service.EstablishToken(userID, code); err != nil {
		fmt.Println(err)
		http.Error(w, "Establish token failed", http.StatusBadRequest)
	}

	// http.Redirect(w, r, "/", http.StatusFound)
	http.Redirect(w, r, "http://localhost:4321/?spotify=connected", http.StatusTemporaryRedirect)
}
