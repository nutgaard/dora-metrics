package dora_metrics_api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"nutgaard/dora-metrics/internal/config"
	"nutgaard/dora-metrics/internal/repositories"
	"nutgaard/dora-metrics/internal/routers"
	"nutgaard/dora-metrics/internal/utils"
	"time"
)

func RunApp(config *config.Config) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Info().Msgf("Loaded config: %s", config)

	deploymentRepository := repositories.CreateDeploymentRepository(config.DB)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		utils.WriteText(writer, "Dora-metrics")
	})
	router.Mount("/api/internal", routers.CreateStatusRouter())
	router.Mount("/api/deployment", routers.CreateDeploymentRouter(deploymentRepository))
	router.Mount("/api", routers.CreateJsonRouter())

	http.ListenAndServe(":"+config.Port, router)
}
