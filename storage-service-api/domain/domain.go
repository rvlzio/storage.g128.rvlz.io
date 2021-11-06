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

type IDConstructor struct{}

func (IDConstructor) NewWarehouseID() WarehouseID {
	return WarehouseID{id: uuid.NewString()}
}

func (IDConstructor) NewFileID() FileID {
	return FileID{id: uuid.NewString()}
}

func (IDConstructor) NewWarehouseStorageID() WarehouseStorageID {
	return WarehouseStorageID{id: uuid.NewString()}
}

func (IDConstructor) NewWarehouseIDFromStr(id string) WarehouseID {
	return WarehouseID{id: id}
}

func (IDConstructor) NewFileIDFromStr(id string) FileID {
	return FileID{id: id}
}

func (IDConstructor) NewWarehouseStorageIDFromStr(id string) WarehouseStorageID {
	return WarehouseStorageID{id: id}
}
