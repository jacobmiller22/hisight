package main

import (
	"fmt"
	"net/http"
)

func (app *hsApp) routes() http.Handler {

	mux := http.DefaultServeMux

	mux.HandleFunc("GET /commands", app.getCommands)
	//
	// mux.HandleFunc("/api/apps", app.applications)
	//
	// mux.HandleFunc("/api/apps/{clientId}", app.application)
	return mux
}

func (app *hsApp) serveHttp() error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.http_port),
		Handler: app.routes(),
	}

	app.log.Info("Starting HTTP Server", "addr", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		app.log.Error("Error starting HTTP Server", "err", err)
		return err
	}
	app.log.Info("HTTP Server Listening", "addr", server.Addr)

	return nil
}
