package models

import (
	"github.com/segmentio/ksuid"
	"time"
)

type Deployment struct {
	Id            ksuid.KSUID
	CreatedAt     time.Time
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
