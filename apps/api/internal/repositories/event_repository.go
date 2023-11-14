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

type event struct {
	Id                  string
	EventType           string
	CreatedAt           time.Time
	Status              string
	RepositoryUrl       string
	Environment         string
	Metadata            string
	OpenedAt            *time.Time
	ClosedAt            *time.Time
	DeploymentReference *string
}

func (d *event) dto() (*models.Event, error) {
	id, err := ksuid.Parse(d.Id)
	if err != nil {
		return nil, err
	}

	eventType, err := models.EventTypeFromString(d.EventType)
	if err != nil {
		return nil, err
	}

	status, err := models.StatusFromString(d.Status)
	if err != nil {
		return nil, err
	}

	var metadata models.Metadata
	err = json.Unmarshal([]byte(d.Metadata), &metadata)
	if err != nil {
		return nil, err
	}

	var deploymentReference *ksuid.KSUID = nil
	if d.DeploymentReference != nil {
		reference, err := ksuid.Parse(*d.DeploymentReference)
		if err != nil {
			return nil, err
		}
		deploymentReference = &reference
	}

	out := &models.Event{
		Id:                  id,
		EventType:           eventType,
		CreatedAt:           d.CreatedAt,
		Status:              status,
		RepositoryUrl:       d.RepositoryUrl,
		Environment:         d.Environment,
		Metadata:            metadata,
		OpenedAt:            d.OpenedAt,
		ClosedAt:            d.ClosedAt,
		DeploymentReference: deploymentReference,
	}

	return out, nil
}

type EventRepositoryCtx struct {
	pool *pgxpool.Pool
}

func (db *EventRepositoryCtx) Ping() error {
	ctx, cancel := DefaultDbContext()
	defer cancel()

	var result int
	err := pgxscan.Get(ctx, db.pool, &result, `SELECT 1`)

	return err
}

func (db *EventRepositoryCtx) GetAll() ([]*models.Event, error) {
	ctx, cancel := DefaultDbContext()
	defer cancel()

	var result []*event
	err := pgxscan.Select(ctx, db.pool, &result, `SELECT * FROM event`)

	out := utils.Map(result, func(element *event, i int) *models.Event {
		v, e := element.dto()
		if e != nil {
			err = e
		}
		return v
	})

	return out, err
}

func (db *EventRepositoryCtx) GetById(id *ksuid.KSUID) (*models.Event, error) {
	ctx, cancel := DefaultDbContext()
	defer cancel()

	var result event
	err := pgxscan.Get(ctx, db.pool, &result, `SELECT * FROM event WHERE id = $1`, id.String())
	if err != nil {
		return nil, err
	}

	return result.dto()
}

func (db *EventRepositoryCtx) GetByType(eventType models.EventType) ([]*models.Event, error) {
	ctx, cancel := DefaultDbContext()
	defer cancel()

	var result []*event
	err := pgxscan.Select(ctx, db.pool, &result, `SELECT * FROM event WHERE event_type = $1`, eventType)

	out := utils.Map(result, func(element *event, i int) *models.Event {
		v, e := element.dto()
		if e != nil {
			err = e
		}
		return v
	})

	return out, err
}

func (db *EventRepositoryCtx) Create(data *models.CreateEvent) (*models.Event, error) {
	panic("implement me")
}

func (db *EventRepositoryCtx) Delete(id *ksuid.KSUID) error {
	ctx, cancel := DefaultDbContext()
	defer cancel()

	_, err := db.pool.Exec(ctx, `DELETE FROM event where id = $1`, id.String())

	return err
}

type EventRepository interface {
	Repository[models.Event, models.CreateEvent]
	GetByType(eventType models.EventType) ([]*models.Event, error)
}

func CreateEventRepository(pool *pgxpool.Pool) EventRepository {
	return &EventRepositoryCtx{pool}
}
