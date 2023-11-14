package models

import (
	"fmt"
	"github.com/segmentio/ksuid"
	"time"
)

type EventType string
type Status string

const (
	PullRequest EventType = "pullrequest"
	Incident    EventType = "incident"
)

var eventtypeMap = map[EventType]struct{}{
	PullRequest: {},
	Incident:    {},
}

func EventTypeFromString(str string) (EventType, error) {
	value := EventType(str)
	_, ok := eventtypeMap[value]
	if !ok {
		return value, fmt.Errorf(`cannot parse: [%s] into EventType`, str)
	}
	return value, nil
}

const (
	Open      Status = "open"
	Closed    Status = "closed"
	Merged    Status = "merged"
	Cancelled Status = "cancelled"
)

var statusMap = map[Status]struct{}{
	Open:      {},
	Closed:    {},
	Merged:    {},
	Cancelled: {},
}

func StatusFromString(str string) (Status, error) {
	value := Status(str)
	_, ok := statusMap[value]
	if !ok {
		return value, fmt.Errorf(`cannot parse: [%s] into Status`, str)
	}
	return value, nil
}

type Event struct {
	Id                  ksuid.KSUID
	EventType           EventType
	CreatedAt           time.Time
	Status              Status
	RepositoryUrl       string
	Environment         string
	Metadata            Metadata
	OpenedAt            *time.Time
	ClosedAt            *time.Time
	DeploymentReference *ksuid.KSUID
}

type CreateEvent struct {
}

func (d *CreateEvent) Validate() error {
	panic("implement me")
}
