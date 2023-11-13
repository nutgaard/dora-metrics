package routers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/segmentio/ksuid"
	renderPkg "github.com/unrolled/render"
	"net/http"
	"nutgaard/dora-metrics/internal/repositories"
	"strconv"
)

var render = renderPkg.New()

func CreateDeploymentRouter(repository *repositories.DeploymentRepository) http.Handler {
	router := chi.NewRouter()

	router.Get("/ping", func(writer http.ResponseWriter, request *http.Request) {
		result, err := repository.Ping()

		respond(map[string]string{"alive": strconv.Itoa(result)}, err, writer)
	})

	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		all, err := repository.GetAll()
		respond(all, err, writer)
	})

	router.Get("/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id, err := ksuid.Parse(chi.URLParam(request, "id"))
		if err != nil {
			respondWithError(err, writer)
			return
		}
		all, err := repository.GetById(id)
		respond(all, err, writer)
	})

	router.Post("/", func(writer http.ResponseWriter, request *http.Request) {
		var d repositories.CreateDeploymentRequest
		err := json.NewDecoder(request.Body).Decode(&d)
		if err != nil {
			respondWithError(err, writer)
			return
		}
		validationErr := d.Validate()
		if validationErr != nil {
			respondWithError(validationErr, writer)
			return
		}

		stored, err := repository.Create(&d)
		if err != nil {
			respondWithError(err, writer)
			return
		}
		respond(stored, err, writer)
	})

	return router
}

func respond[A any](value A, err error, writer http.ResponseWriter) {
	if err != nil {
		respondWithError(err, writer)
		return
	}
	err = render.JSON(writer, 200, value)
	if err != nil {
		respondWithError(err, writer)
		return
	}
}

func respondWithError(err error, writer http.ResponseWriter) {
	writer.WriteHeader(500)
	_, _ = writer.Write([]byte(err.Error()))
}
