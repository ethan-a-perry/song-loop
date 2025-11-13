package main

import (
	"fmt"
	"net/http"

	"github.com/ethan-a-perry/song-loop/internal/auth"
)

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	err := r.URL.Query().Get("error")

	if err != "" {
		http.Error(w, "Authorization failed during callback" + err, http.StatusUnauthorized)
		return
	}

	code := r.URL.Query().Get("code")

	if code == "" {
		http.Error(w, "No code provided", http.StatusBadRequest)
        return
	}

	auth.EstablishToken(code)

	http.Redirect(w, r, "/", http.StatusFound)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	authorized, err := auth.Authenticate()

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if authorized {
		return
	}

	// Not authorized
	authorizationUrl, err := auth.GetAuthorizationUrl()

	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, authorizationUrl, http.StatusFound)
}

func main() {
	if err := auth.InitSpotify(); err != nil {
		return
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/callback", callbackHandler)

	fmt.Println("Server running at http://127.0.0.1:8080")
	http.ListenAndServe(":8080", nil)
}
