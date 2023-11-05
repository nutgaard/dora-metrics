package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/ksuid"
	ConfigLoader "nutgaard/dora-metrics/internal/config"
	"nutgaard/dora-metrics/internal/models"
	"nutgaard/dora-metrics/internal/utils"
	"time"
)

type DeploymentRepository struct {
	pool *pgxpool.Pool
}

type deployment struct {
	Id            string
	CreatedAt     time.Time `db:"created_at"`
	StartedAt     time.Time `db:"started_at"`
	FinishedAt    time.Time `db:"finished_at"`
	RepositoryUrl string    `db:"repository_url"`
	Application   string
	Environment   string
	Department    *string
	Team          *string
	Product       *string
	Version       *string
}

func (d *deployment) dto() (*models.Deployment, error) {
	id, err := ksuid.Parse(d.Id)
	if err != nil {
		return nil, err
	}

	out := &models.Deployment{
		Id:            id,
		CreatedAt:     d.CreatedAt,
		StartedAt:     d.StartedAt,
		FinishedAt:    d.FinishedAt,
		RepositoryUrl: d.RepositoryUrl,
		Application:   d.Application,
		Environment:   d.Environment,
		Department:    d.Department,
		Team:          d.Team,
		Product:       d.Product,
		Version:       d.Version,
	}

	return out, nil
}

type CreateDeploymentRequest struct {
	StartedAt     time.Time
	FinishedAt    time.Time
	RepositoryUrl string
	Application   string
	Environment   string
	Department    *string
	Team          *string
	Product       *string
	Version       *string
}

func (d *CreateDeploymentRequest) Validate() error {
	if d.StartedAt.IsZero() {
		return errors.New(`'StartedAt' cannot be empty`)
	}
	if d.FinishedAt.IsZero() {
		return errors.New(`'FinishedAt' cannot be empty`)
	}
	if len(d.RepositoryUrl) == 0 {
		return errors.New(`'RepositoryUrl' cannot be empty`)
	}
	if len(d.Application) == 0 {
		return errors.New(`'Application' cannot be empty`)
	}
	if len(d.Environment) == 0 {
		return errors.New(`'Environment' cannot be empty`)
	}
	return nil
}

func CreateDeploymentRepository(config ConfigLoader.PostgresqlConfig) *DeploymentRepository {
	connectionUrl := fmt.Sprintf("postgresql://%s:%s@%s", config.Username, config.Password, config.Host)
	pool, err := pgxpool.New(context.Background(), connectionUrl)
	if err != nil {
		log.Fatal().AnErr("error", err).Msg("Could not connect to db")
	}

	return &DeploymentRepository{pool}
}

func (db DeploymentRepository) Ping() (int, error) {
	ctx, cancel := defaultContext()
	defer cancel()

	var result int
	err := pgxscan.Get(ctx, db.pool, &result, `SELECT 101`)

	return result, err
}

func (db DeploymentRepository) GetAll() ([]*models.Deployment, error) {
	ctx, cancel := defaultContext()
	defer cancel()

	var result []*deployment
	err := pgxscan.Select(ctx, db.pool, &result, `SELECT * FROM deployment`)

	out := utils.Map(result, func(element *deployment, i int) *models.Deployment {
		v, e := element.dto()
		if e != nil {
			err = e
		}
		return v
	})

	return out, err
}

func (db DeploymentRepository) GetById(id ksuid.KSUID) (*models.Deployment, error) {
	ctx, cancel := defaultContext()
	defer cancel()

	var result deployment
	err := pgxscan.Get(ctx, db.pool, &result, `SELECT * FROM deployment WHERE id = $1`, id.String())
	if err != nil {
		return nil, err
	}

	out, err := result.dto()
	return out, err
}

func (db DeploymentRepository) Create(d *CreateDeploymentRequest) (*models.Deployment, error) {
	deployment := models.Deployment{
		Id:            ksuid.New(),
		CreatedAt:     time.Now(),
		StartedAt:     d.StartedAt,
		FinishedAt:    d.FinishedAt,
		RepositoryUrl: d.RepositoryUrl,
		Application:   d.Application,
		Environment:   d.Environment,
		Department:    d.Department,
		Team:          d.Team,
		Product:       d.Product,
		Version:       d.Version,
	}

	return &deployment, nil
}

func defaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
