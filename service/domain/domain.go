package domain

import (
	"github.com/google/uuid"
)

type Message interface {
	IsMessage()
}

type Event interface {
	IsEvent()
}

type Aggregate interface {
	Events() []Event
}

type WarehouseID struct {
	id string
}

func (w WarehouseID) Str() string {
	return w.id
}

type FileID struct {
	id string
}

func (f FileID) Str() string {
	return f.id
}

type WarehouseStorageID struct {
	id string
}

func (wf WarehouseStorageID) Str() string {
	return wf.id
}

type IDFactory struct{}

func (factory IDFactory) NewWarehouseID() WarehouseID {
	return WarehouseID{id: uuid.NewString()}
}

func (factory IDFactory) NewFileID() FileID {
	return FileID{id: uuid.NewString()}
}

func (factory IDFactory) NewWarehouseStorageID() WarehouseStorageID {
	return WarehouseStorageID{id: uuid.NewString()}
}

func (factory IDFactory) NewWarehouseIDFromStr(id string) WarehouseID {
	return WarehouseID{id: id}
}

func (factory IDFactory) NewFileIDFromStr(id string) FileID {
	return FileID{id: id}
}

func (factory IDFactory) NewWarehouseStorageIDFromStr(id string) WarehouseStorageID {
	return WarehouseStorageID{id: id}
}
