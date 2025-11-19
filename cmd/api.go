package main

import (
	"fmt"
	"net/http"

	"github.com/ethan-a-perry/song-loop/internal/auth"
	"github.com/ethan-a-perry/song-loop/internal/database/data"
	"github.com/ethan-a-perry/song-loop/internal/database/dataaccess"
	"github.com/ethan-a-perry/song-loop/internal/middleware"

	// "github.com/ethan-a-perry/song-loop/internal/middleware"
	// "github.com/ethan-a-perry/song-loop/internal/spotify"
	"github.com/ethan-a-perry/song-loop/internal/spotifyauth"
)

type api struct {
	config config
	db *dataaccess.MongoDataAccess
}

type config struct {
	addr string
}

func (a *api) mount() http.Handler {
	router := http.NewServeMux()

	userData := data.NewUserData(a.db.UserCollection)

	// Auth
	authService := auth.NewService(userData)
	authHandler := auth.NewHandler(authService)

	protect := func(handler http.HandlerFunc) http.Handler {
		return middleware.CORS(authService.Middleware(http.HandlerFunc(handler)))
	}

	router.Handle("/api/user", protect(authHandler.GetUserFromRequest))

	// Spotify Auth
	spotifyAuthService := spotifyauth.NewService(userData)
	spotifyAuthHandler := spotifyauth.NewHandler(spotifyAuthService)

	router.Handle("/api/spotify/connect", protect(spotifyAuthHandler.Connect))
	router.HandleFunc("/api/spotify/callback", spotifyAuthHandler.Callback)

	// Spotify
	// spotifyService := spotify.NewService()
	// spotifyHandler := spotify.NewHandler(spotifyService)

	// router.HandleFunc("/loop", spotifyHandler.Loop)

	return router
}

func (a *api) run(router http.Handler) error {
	server := http.Server {
		Addr: a.config.addr,
		Handler: router,
	}

	fmt.Println("Server running at http://127.0.0.1:8080")

	return server.ListenAndServe()
}
