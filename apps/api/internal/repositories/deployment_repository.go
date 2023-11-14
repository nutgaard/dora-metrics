package repositories

import (
	"encoding/json"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/ksuid"
	"nutgaard/dora-metrics/internal/models"
	"nutgaard/dora-metrics/internal/utils"
	"time"
)

type deployment struct {
	Id            string
	CreatedAt     time.Time
	StartedAt     time.Time
	FinishedAt    time.Time
	RepositoryUrl string
	Environment   string
	Metadata      string
}

func (d *deployment) dto() (*models.Deployment, error) {
	id, err := ksuid.Parse(d.Id)
	if err != nil {
		return nil, err
	}

	var metadata models.Metadata
	err = json.Unmarshal([]byte(d.Metadata), &metadata)
	if err != nil {
		return nil, err
	}

	out := &models.Deployment{
		Id:            id,
		CreatedAt:     d.CreatedAt,
		StartedAt:     d.StartedAt,
		FinishedAt:    d.FinishedAt,
		RepositoryUrl: d.RepositoryUrl,
		Environment:   d.Environment,
		Metadata:      metadata,
	}

	return out, nil
}

type DeploymentRepositoryCtx struct {
	pool *pgxpool.Pool
}

func (db *DeploymentRepositoryCtx) Ping() error {
	ctx, cancel := DefaultDbContext()
	defer cancel()

	var result int
	err := pgxscan.Get(ctx, db.pool, &result, `SELECT 1`)

	return err
}

func (db *DeploymentRepositoryCtx) GetAll() ([]*models.Deployment, error) {
	ctx, cancel := DefaultDbContext()
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

func (db *DeploymentRepositoryCtx) GetById(id *ksuid.KSUID) (*models.Deployment, error) {
	ctx, cancel := DefaultDbContext()
	defer cancel()

	var result deployment
	err := pgxscan.Get(ctx, db.pool, &result, `SELECT * FROM deployment WHERE id = $1`, id.String())
	if err != nil {
		return nil, err
	}

	return result.dto()
}

func (db *DeploymentRepositoryCtx) Create(data *models.CreateDeployment) (*models.Deployment, error) {
	panic("implement me")
}

func (db *DeploymentRepositoryCtx) Delete(id *ksuid.KSUID) error {
	ctx, cancel := DefaultDbContext()
	defer cancel()

	_, err := db.pool.Exec(ctx, `DELETE FROM deployment where id = $1`, id.String())

	return err
}

type DeploymentRepository interface {
	Repository[models.Deployment, models.CreateDeployment]
}

func CreateDeploymentRepository(pool *pgxpool.Pool) DeploymentRepository {
	return &DeploymentRepositoryCtx{pool}
}
