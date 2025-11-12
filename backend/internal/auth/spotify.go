package auth

import (
	"errors"
	"fmt"
	"net/url"
)

func GetAuthorizationUrl(clientId, redirectUri, scope string) (string, error) {
	codeVerifier, err := GenerateCodeVerifier()

	if err != nil {
		return "", errors.New("Could not generate code verifier")
	}

	codeChallenge := GenerateCodeChallenge(codeVerifier)

	state, err := GenerateCodeVerifier()
	if err != nil {
		return "", errors.New("Could not generate state")
	}

	v := url.Values{}
	v.Set("client_id", clientId)
	v.Set("response_type", "code")
	v.Set("redirect_uri", redirectUri)
	v.Set("state", state)
	v.Set("scope", scope)
	v.Set("code_challenge_method", "S256")
	v.Set("code_challenge", codeChallenge)

	return fmt.Sprintf("https://accounts.spotify.com/authorize?%s", v.Encode()), nil
}
