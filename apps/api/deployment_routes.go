package main

import (
	"github.com/go-chi/chi/v5"
	renderPkg "github.com/unrolled/render"
	"net/http"
	"strconv"
)

func deploymentRouter(repository *DeploymentRepository) http.Handler {
	render := renderPkg.New()
	router := chi.NewRouter()

	router.Get("/ping", func(writer http.ResponseWriter, request *http.Request) {
		result := repository.Ping()

		err := render.JSON(writer, 200, map[string]string{"alive": strconv.Itoa(result)})
		if err != nil {
			writer.WriteHeader(500)
			return
		}
	})

	return router
}
