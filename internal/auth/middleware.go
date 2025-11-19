package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func (s *svc) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[AUTH] Method: '%s' Path: '%s'\n", r.Method, r.URL.Path)
		fmt.Printf("[AUTH] Full URL: %+v\n", r.URL)

		fmt.Println("receieved auth request")
		keyset, err := jwk.Fetch(r.Context(), "http://localhost:4321/api/auth/jwks")
		if err != nil {
			fmt.Println("1")
			http.Error(w, "failed to fetch jwks", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseRequest(r, jwt.WithKeySet(keyset))
		if err != nil {
			fmt.Println("2")
			http.Error(w, "failed to parse request", http.StatusUnauthorized)
			return
		}

		_, exists := token.Subject()
		if !exists {
			fmt.Println("3")
			http.Error(w, "failed to find user id in token", http.StatusUnauthorized)
			return
		}

		fmt.Println("4")

		var id string
		var email string

		token.Get("id", &id)
		token.Get("email", &email)

		ctx := context.WithValue(r.Context(), "user_id", id)
		ctx = context.WithValue(ctx, "email", email)

		// next.ServeHTTP(w, r.WithContext(ctx))

		fmt.Println("[AUTH] Calling next.ServeHTTP")
				next.ServeHTTP(w, r.WithContext(ctx))
				fmt.Println("[AUTH] Returned from next.ServeHTTP")
	})
}

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("user_id").(string)

	return userID, ok
}

func GetUserEmail(ctx context.Context) (string, bool) {
	email, ok := ctx.Value("email").(string)

	return email, ok
}
