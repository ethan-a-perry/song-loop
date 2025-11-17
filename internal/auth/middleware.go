package auth

import (
	"context"
	"net/http"

	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func (s *svc) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keyset, err := jwk.Fetch(r.Context(), "http://localhost:4321/api/auth/jwks")
		if err != nil {
			http.Error(w, "failed to fetch jwks", http.StatusUnauthorized)
		}

		token, err := jwt.ParseRequest(r, jwt.WithKeySet(keyset))
		if err != nil {
			http.Error(w, "failed to parse request", http.StatusUnauthorized)
		}

		_, exists := token.Subject()
		if !exists {
			http.Error(w, "failed to find user id in token", http.StatusUnauthorized)
		}

		var id string
		var email string

		token.Get("id", &id)
		token.Get("email", &email)

		ctx := context.WithValue(r.Context(), "user_id", id)
		ctx = context.WithValue(ctx, "email", email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetContextValue(r *http.Request, value string) string {
	userID, _ := r.Context().Value(value).(string)

	return userID
}
