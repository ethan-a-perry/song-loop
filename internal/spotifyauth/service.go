package spotifyauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Service interface {
	GetAuthorizationUrl(state string) (string, error)
	EstablishToken(userID, code string) error
	Authenticate(userID string) (*models.SpotifyToken, error)
}

type svc struct {
	userData *data.UserData
}

var codeVerifier string

func NewService(userData *data.UserData) Service {
	return &svc {
		userData: userData,
	}
}

func (s *svc) GetAuthorizationUrl(state string) (string, error) {
	var err error
	codeVerifier, err = GenerateCodeVerifier(64)
	if err != nil {
		return "", errors.New("Could not generate code verifier")
	}

	codeChallenge := GenerateCodeChallenge(codeVerifier)

	v := url.Values{}
	v.Set("client_id", os.Getenv("CLIENT_ID"))
	v.Set("response_type", "code")
	v.Set("redirect_uri", os.Getenv("REDIRECT_URI"))
	v.Set("state", state)
	v.Set("scope", os.Getenv("SCOPE"))
	v.Set("code_challenge_method", "S256")
	v.Set("code_challenge", codeChallenge)

	return fmt.Sprintf("https://accounts.spotify.com/authorize?%s", v.Encode()), nil
}

func (s *svc) GetToken(userID string, data url.Values) error {
	req, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))

	if err != nil {
		return fmt.Errorf("Request creation failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("Request from client failed: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("spotify token exchange failed: %s", body)
	}

	var tr struct {
		AccessToken string `json:"access_token"`
		TokenType string `json:"token_type"`
		Scope string `json:"scope"`
		ExpiresIn int `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(res.Body).Decode(&tr); err != nil {
		return fmt.Errorf("failed to decode token response: %v", err)
	}

	token := models.SpotifyToken{
		AccessToken: tr.AccessToken,
		TokenType: tr.TokenType,
		Scope: tr.Scope,
		ExpiresAt: time.Now().Add((time.Duration(tr.ExpiresIn) - 10) * time.Second),
		RefreshToken: tr.RefreshToken,
	}

	if err := s.SaveToken(userID, &token); err != nil {
        return fmt.Errorf("Failed to save token: %w", err)
    }

    fmt.Println(token)

	return nil
}

func (s *svc) EstablishToken(userID, code string) error {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", os.Getenv("REDIRECT_URI"))
	data.Set("client_id", os.Getenv("CLIENT_ID"))
	data.Set("code_verifier", codeVerifier)

	return s.GetToken(userID, data)
}

func (s *svc) RefreshToken(userID string, refreshToken string) error {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", os.Getenv("CLIENT_ID"))

	return s.GetToken(userID, data)
}

func (s *svc) SaveToken(userID string, token *models.SpotifyToken) error {
	if err := s.userData.UpdateSpotifyToken(userID, token); err != nil {
        return fmt.Errorf("Failed to save token: %w", err)
    }

    return nil
}

func (s *svc) LoadToken(userID string) (*models.SpotifyToken, error) {
	user, err := s.userData.GetUserById(userID)
	if err != nil {
		return &models.SpotifyToken{}, fmt.Errorf("Could not find user by the id in database: %w", err)
	}

	if user.SpotifyToken == (models.SpotifyToken{}) {
		return &models.SpotifyToken{}, fmt.Errorf("spotify token is missing: %w", err)
	}

	return &user.SpotifyToken, nil
}

func (s *svc) Authenticate(userID string) (*models.SpotifyToken, error) {
	token, err := s.LoadToken(userID)
	if err != nil {
		return &models.SpotifyToken{}, fmt.Errorf("failed to find spotify token for user: %w", err)
	}

	if time.Now().After(token.ExpiresAt) {
		err := s.RefreshToken(userID, token.RefreshToken)
		if err != nil {
			return &models.SpotifyToken{}, fmt.Errorf("Failed to refresh expired token: %w", err)
		}

		token, err = s.LoadToken(userID)
		if err != nil {
			return &models.SpotifyToken{}, fmt.Errorf("failed to find spotify token for user: %w", err)
		}
	}

	return token, nil
}
