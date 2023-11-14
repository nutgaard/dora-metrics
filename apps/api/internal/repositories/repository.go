package repositories

import "github.com/segmentio/ksuid"

type Repository[Type any, CreateType any] interface {
	Ping() error
	GetAll() ([]*Type, error)
	GetById(id *ksuid.KSUID) (*Type, error)
	Create(data *CreateType) (*Type, error)
	Delete(id *ksuid.KSUID) error
}
