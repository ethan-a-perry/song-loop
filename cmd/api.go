package main

import (
	"fmt"
	"net/http"

	"github.com/ethan-a-perry/song-loop/internal/spotify"
	"github.com/ethan-a-perry/song-loop/internal/spotifyauth"
	"github.com/ethan-a-perry/song-loop/internal/store"
	"github.com/ethan-a-perry/song-loop/internal/web"
)

type api struct {
	config config
	store *store.Store
	authService *spotifyauth.Service
    spotifyService *spotify.Service
    webService *web.Service
}

type config struct {
	addr string
}

func (a *api) mount() http.Handler {
	router := http.NewServeMux()

	// Auth
	authHandler := spotifyauth.NewHandler(a.authService)
	router.HandleFunc("/api/spotify/connect", authHandler.Connect)
	router.HandleFunc("/api/spotify/callback", authHandler.Callback)

	// Loop
	spotifyHandler := spotify.NewHandler(a.spotifyService)
	router.HandleFunc("POST /api/spotify/loop", spotifyHandler.Loop)
	router.HandleFunc("/api/spotify/loop/stop", spotifyHandler.StopLoop)

	// Static files
	router.Handle("/web/static/", http.StripPrefix("/web/static/", http.FileServer(http.Dir("web/static"))))

	// Web
	webHandler := web.NewHandler(a.webService)
	router.HandleFunc("/", webHandler.Index)

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
