package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	ConfigLoader "nutgaard/dora-metrics/config"
	Utils "nutgaard/dora-metrics/utils"
	"time"
)

func runApp(config *ConfigLoader.Config) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Info().Msgf("Loaded config: %s", config)

	deploymentRepository := CreateDeploymentRepository(config.DB)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		Utils.WriteText(writer, "Dora-metrics")
	})
	router.Mount("/api/internal", statusRouter())
	router.Mount("/api/deployment", deploymentRouter(deploymentRepository))
	router.Mount("/api", jsonRouter())

	http.ListenAndServe(":"+config.Port, router)
}
