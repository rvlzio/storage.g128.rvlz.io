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

func (queue *Queue) IsWaiting(file File) bool {
	for _, reservation := range queue.reservations {
		if reservation.File.ID == file.ID {
			return true
		}
	}
	return false
}

func (queue *Queue) Add(file File) {
	reservation := Reservation{File: file}
	queue.reservations = append(queue.reservations, reservation)
}

func (queue *Queue) Remove(fileID dm.FileID) (File, error) {
	for ix, reservation := range queue.reservations {
		if reservation.File.ID == fileID {
			if ix == 0 {
				queue.reservations = queue.reservations[1:]
			} else if n := len(queue.reservations); ix == n-1 {
				queue.reservations = queue.reservations[:ix-1]
			} else {
				queue.reservations = append(
					queue.reservations[:ix-1],
					queue.reservations[ix+1:]...,
				)
			}
			return reservation.File, nil
		}
	}
	return File{}, nil
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

func (ws *WarehouseStorage) WarehouseID() dm.WarehouseID {
	return ws.warehouseID
}

func (ws *WarehouseStorage) AvailableStorage() int {
	return ws.unclaimedStorage - ws.queue.Size()
}

func (ws *WarehouseStorage) ReservedStorage() int {
	return ws.queue.Size()
}

func (ws *WarehouseStorage) Reserve(file File) error {
	if ws.queue.IsWaiting(file) {
		ws.events = append(ws.events, ev.StorageReservationDuplicated{
			WarehouseID: ws.warehouseID,
			FileID:      file.ID,
		})
		return er.StorageReservationDuplicated
	}
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

func (ws *WarehouseStorage) Unreserve(fileID dm.FileID) error {
	if !ws.queue.IsWaiting(File{fileID, 0}) {
		ws.events = append(ws.events, ev.NonexistentStorageReservationUnreserved{
			WarehouseID: ws.warehouseID,
			FileID:      fileID,
		})
		return er.NonexistentStorageReservationUnreserved
	}
	file, _ := ws.queue.Remove(fileID)
	ws.events = append(ws.events, ev.StorageUnreserved{
		WarehouseID:       ws.warehouseID,
		FileID:            file.ID,
		UnreservedStorage: file.Size,
	})
	return nil
}

func (ws *WarehouseStorage) Commit(fileID dm.FileID) error {
	if !ws.queue.IsWaiting(File{fileID, 0}) {
		ws.events = append(ws.events, ev.UnreservedStorageCommitted{
			WarehouseID: ws.warehouseID,
			FileID:      fileID,
		})
		return er.UnreservedStorageCommitted
	}
	file, _ := ws.queue.Remove(fileID)
	ws.unclaimedStorage -= file.Size
	ws.events = append(ws.events, ev.ReservedStorageCommitted{
		WarehouseID: ws.warehouseID,
		FileID:      fileID,
	})
	return nil
}

func (ws *WarehouseStorage) Free(file File) error {
	if ws.queue.IsWaiting(file) {
		ws.events = append(ws.events, ev.FreeingUncommittedStorageAttempted{
			WarehouseID: ws.WarehouseID(),
			FileID:      file.ID,
		})
		return er.FreeingUncommittedStorageAttempted
	}
	if file.Size > ws.AvailableStorage() {
		ws.events = append(ws.events, ev.FreedStorageExceededAvailability{
			WarehouseID:      ws.WarehouseID(),
			FileID:           file.ID,
			AvailableStorage: ws.AvailableStorage(),
		})
		return er.FreedStorageExceededAvailability
	}
	ws.unclaimedStorage += file.Size
	ws.events = append(ws.events, ev.StorageFreed{
		WarehouseID:  ws.warehouseID,
		FileID:       file.ID,
		FreedStorage: file.Size,
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
