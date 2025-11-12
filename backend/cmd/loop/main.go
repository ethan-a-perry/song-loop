package main

import (
	"fmt"
	"net/http"

	"github.com/ethan-a-perry/loop/internal/auth"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	authorizationUrl, err := auth.GetAuthorizationUrl(
		"",
		"http://127.0.0.1:8080",
		"user-read-private user-read-email",
	)

	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, authorizationUrl, http.StatusFound)
}

func main() {
	http.HandleFunc("/login", loginHandler)

	fmt.Println("Server running at http://127.0.0.1:8080")
	http.ListenAndServe(":8080", nil)
}
