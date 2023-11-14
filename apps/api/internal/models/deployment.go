package models

import (
	"github.com/segmentio/ksuid"
	"time"
)

type Deployment struct {
	Id         ksuid.KSUID
	CreatedAt  time.Time
	StartedAt  time.Time
	FinishedAt time.Time

	RepositoryUrl string
	Environment   string
	Metadata      Metadata
}

type CreateDeployment struct {
}

func (d *CreateDeployment) Validate() error {
	panic("implement me")
}
