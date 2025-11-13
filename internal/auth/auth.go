package auth

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

	"github.com/ethan-a-perry/song-loop/config"
)

var (
	session *config.Config
	codeVerifier string
)

type Token struct {
	AccessToken string
	TokenType string
	Scope string
	ExpiresAt time.Time
	RefreshToken string
}

func InitSpotify() error {
	var err error
	session, err = config.LoadConfig()

	if err != nil {
		return fmt.Errorf("Failed to load config: %w", err)
	}

	return nil
}

func GetAuthorizationUrl() (string, error) {
	var err error
	codeVerifier, err = GenerateCodeVerifier()
	if err != nil {
		return "", errors.New("Could not generate code verifier")
	}

	codeChallenge := GenerateCodeChallenge(codeVerifier)

	state, err := GenerateCodeVerifier()
	if err != nil {
		return "", errors.New("Could not generate state")
	}

	v := url.Values{}
	v.Set("client_id", session.ClientID)
	v.Set("response_type", "code")
	v.Set("redirect_uri", session.RedirectURI)
	v.Set("state", state)
	v.Set("scope", session.Scope)
	v.Set("code_challenge_method", "S256")
	v.Set("code_challenge", codeChallenge)

	return fmt.Sprintf("https://accounts.spotify.com/authorize?%s", v.Encode()), nil
}

func GetToken(data url.Values) error {
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

	token := Token{
		AccessToken: tr.AccessToken,
		TokenType: tr.TokenType,
		Scope: tr.Scope,
		ExpiresAt: time.Now().Add((time.Duration(tr.ExpiresIn) - 10) * time.Second),
		RefreshToken: tr.RefreshToken,
	}

	if err := SaveToken(&token); err != nil {
        return fmt.Errorf("Failed to save token: %w", err)
    }

	return nil
}

func EstablishToken(code string) error {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", session.RedirectURI)
	data.Set("client_id", session.ClientID)
	data.Set("code_verifier", codeVerifier)

	return GetToken(data)
}

func RefreshToken(refreshToken string) error {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", session.ClientID)

	return GetToken(data)
}

func LoadToken() (*Token, error) {
	data, err := os.ReadFile("internal/auth/token.json")

	if err != nil {
		return nil, fmt.Errorf("Failed to read token file: %w", err)
	}

	var t Token
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("Failed to parse token file: %w", err)
	}

	return &t, nil
}

func SaveToken(t *Token) error {
	data, err := json.MarshalIndent(t, "", "  ")

	if err != nil {
		return fmt.Errorf("Failed to marshal token file to JSON: %w", err)
	}

	if err := os.WriteFile("internal/auth/token.json", data, 0600); err != nil {
        return fmt.Errorf("Failed to write token file: %w", err)
    }

    return nil
}

func Authenticate() (*Token, bool, error) {
	t, err := LoadToken()

	if err != nil {
		return t, false, nil
	}

	if time.Now().After(t.ExpiresAt) {
		err := RefreshToken(t.RefreshToken)

		if err != nil {
			return t, false, fmt.Errorf("Failed to refresh token: %w", err)
		}
	}

	return t, true, nil
}
