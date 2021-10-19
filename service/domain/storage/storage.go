package storage

import (
	dm "storage.g128.rvlz.io/domain"
	er "storage.g128.rvlz.io/domain/storage/errors"
	ev "storage.g128.rvlz.io/domain/storage/events"
)

type Reservation struct {
	File
}

type Queue struct {
	reservations []Reservation
}

func (queue *Queue) Add(file File) {
	reservation := Reservation{File: file}
	queue.reservations = append(queue.reservations, reservation)
}

func (queue *Queue) Size() int {
	size := 0
	for _, reservation := range queue.reservations {
		size += reservation.File.Size
	}
	return size
}

type File struct {
	ID   dm.FileID
	Size int
}

type WarehouseStorage struct {
	id               dm.WarehouseStorageID
	warehouseID      dm.WarehouseID
	unclaimedStorage int
	queue            Queue
	events           []dm.Event
}

func (ws *WarehouseStorage) AvailableStorage() int {
	return ws.unclaimedStorage - ws.queue.Size()
}

func (ws *WarehouseStorage) ReservedStorage() int {
	return ws.queue.Size()
}

func (ws *WarehouseStorage) Reserve(file File) error {
	availableStorage := ws.AvailableStorage()
	if file.Size > availableStorage {
		ws.events = append(ws.events, ev.AvailableStorageExceeded{
			WarehouseID:      ws.warehouseID,
			FileID:           file.ID,
			AvailableStorage: availableStorage,
		})
		return er.AvailableStorageExceeded
	}
	ws.queue.Add(file)
	ws.events = append(ws.events, ev.StorageReserved{
		WarehouseID: ws.warehouseID,
		FileID:      file.ID,
	})
	return nil
}

func (ws *WarehouseStorage) Events() []dm.Event {
	return ws.events
}

func (ws *WarehouseStorage) clearEvents() {
	ws.events = nil
}

type StorageFactory struct{}

func (factory StorageFactory) NewWarehouseStorage(
	warehouseID dm.WarehouseID,
	capacity int,
) WarehouseStorage {
	warehouseStorageID := dm.IDFactory{}.NewWarehouseStorageID()
	queue := Queue{reservations: []Reservation{}}
	return WarehouseStorage{
		id:               warehouseStorageID,
		warehouseID:      warehouseID,
		unclaimedStorage: capacity,
		queue:            queue,
	}
}
