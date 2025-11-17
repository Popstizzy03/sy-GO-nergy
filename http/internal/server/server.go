package server

import (
	"net/http"
	"time"
)

func NewServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
