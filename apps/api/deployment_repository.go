package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	ConfigLoader "nutgaard/dora-metrics/config"
	"time"
)

type DeploymentRepository struct {
	pool *pgxpool.Pool
}

func CreateDeploymentRepository(config ConfigLoader.PostgresqlConfig) *DeploymentRepository {
	// postgresql://localhost:8082/dora-metrics
	connectionUrl := fmt.Sprintf("postgresql://%s:%s@%s", config.Username, config.Password, config.Host)
	pool, err := pgxpool.New(context.Background(), connectionUrl)
	if err != nil {
		log.Fatal().AnErr("error", err).Msg("Could not connect to db")
	}

	return &DeploymentRepository{pool}
}

func (db DeploymentRepository) Ping() int {
	ctx, cancel := defaultContext()
	defer cancel()

	var result int
	err := db.pool.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
		return conn.QueryRow(ctx, "select 101").Scan(&result)
	})
	if err != nil {
		log.Fatal().AnErr("error", err).Msg("Error while connecting to db")
	}

	return result
}

func defaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
