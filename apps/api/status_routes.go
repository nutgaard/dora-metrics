package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	Utils "nutgaard/dora-metrics/utils"
)

func statusRouter() http.Handler {
	router := chi.NewRouter()
	router.Get("/isAlive", func(writer http.ResponseWriter, request *http.Request) {
		Utils.WriteText(writer, "Alive")
	})

	router.Get("/isReady", func(writer http.ResponseWriter, request *http.Request) {
		Utils.WriteText(writer, "Ready")
	})

	return router
}
