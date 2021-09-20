package apiserver

import (
	"net/http"
	"pingram/pkg/config"
)

func Start() error {
	srv := newServer()

	return http.ListenAndServe(config.Settings.BindAddr, srv)
}
