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
				queue.reservations = queue.reservations[:ix]
			} else {
				queue.reservations = append(
					queue.reservations[:ix],
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
	id             dm.WarehouseStorageID
	warehouseID    dm.WarehouseID
	capacity       int
	claimedStorage int
	queue          Queue
	events         []dm.Event
}

func (ws *WarehouseStorage) WarehouseID() dm.WarehouseID {
	return ws.warehouseID
}

func (ws *WarehouseStorage) Capacity() int {
	return ws.capacity
}

func (ws *WarehouseStorage) AvailableStorage() int {
	return ws.capacity - ws.claimedStorage - ws.queue.Size()
}

func (ws *WarehouseStorage) ReservedStorage() int {
	return ws.queue.Size()
}

func (ws *WarehouseStorage) ClaimedStorage() int {
	return ws.claimedStorage
}

func (ws *WarehouseStorage) Reserve(file File) error {
	if ws.queue.IsWaiting(file) {
		ws.events = append(ws.events, ev.StorageReservationDuplicated{
			WarehouseID: ws.WarehouseID(),
			FileID:      file.ID,
		})
		return er.StorageReservationDuplicated
	}
	availableStorage := ws.AvailableStorage()
	if file.Size > availableStorage {
		ws.events = append(ws.events, ev.AvailableStorageExceeded{
			WarehouseID:      ws.WarehouseID(),
			FileID:           file.ID,
			AvailableStorage: availableStorage,
		})
		return er.AvailableStorageExceeded
	}
	ws.queue.Add(file)
	ws.events = append(ws.events, ev.StorageReserved{
		WarehouseID: ws.WarehouseID(),
		FileID:      file.ID,
	})
	return nil
}

func (ws *WarehouseStorage) Unreserve(fileID dm.FileID) error {
	if !ws.queue.IsWaiting(File{fileID, 0}) {
		ws.events = append(ws.events, ev.NonexistentStorageReservationUnreserved{
			WarehouseID: ws.WarehouseID(),
			FileID:      fileID,
		})
		return er.NonexistentStorageReservationUnreserved
	}
	file, _ := ws.queue.Remove(fileID)
	ws.events = append(ws.events, ev.StorageUnreserved{
		WarehouseID:       ws.WarehouseID(),
		FileID:            file.ID,
		UnreservedStorage: file.Size,
	})
	return nil
}

func (ws *WarehouseStorage) Commit(fileID dm.FileID) error {
	if !ws.queue.IsWaiting(File{fileID, 0}) {
		ws.events = append(ws.events, ev.UnreservedStorageCommitted{
			WarehouseID: ws.WarehouseID(),
			FileID:      fileID,
		})
		return er.UnreservedStorageCommitted
	}
	file, _ := ws.queue.Remove(fileID)
	ws.claimedStorage += file.Size
	ws.events = append(ws.events, ev.ReservedStorageCommitted{
		WarehouseID: ws.WarehouseID(),
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
	if file.Size > ws.ClaimedStorage() {
		ws.events = append(ws.events, ev.FreedStorageExceededClaimedStorage{
			WarehouseID:    ws.WarehouseID(),
			FileID:         file.ID,
			ClaimedStorage: ws.ClaimedStorage(),
		})
		return er.FreedStorageExceededClaimedStorage
	}
	ws.claimedStorage -= file.Size
	ws.events = append(ws.events, ev.StorageFreed{
		WarehouseID:  ws.WarehouseID(),
		FileID:       file.ID,
		FreedStorage: file.Size,
	})
	return nil
}

func (ws *WarehouseStorage) Expand(capacity int) error {
	if ws.capacity > capacity {
		ws.events = append(ws.events, ev.MinimumStorageExpansionNotMet{
			WarehouseID: ws.WarehouseID(),
			Capacity:    ws.Capacity(),
		})
		return er.MinimumStorageExpansionNotMet
	}
	ws.capacity = capacity
	ws.events = append(ws.events, ev.StorageExpanded{
		WarehouseID:     ws.WarehouseID(),
		ExpandedStorage: capacity,
	})
	return nil
}

func (ws *WarehouseStorage) Shrink(capacity int) error {
	if ws.capacity < capacity {
		ws.events = append(ws.events, ev.MaximumStorageContractionNotMet{
			WarehouseID: ws.WarehouseID(),
			Capacity:    ws.Capacity(),
		})
		return er.MaximumStorageContractionNotMet
	}
	if capacity < ws.Capacity()-ws.AvailableStorage() {
		ws.events = append(ws.events, ev.MinimumStorageContractionNotMet{
			WarehouseID: ws.WarehouseID(),
			Capacity:    ws.Capacity(),
		})
		return er.MinimumStorageContractionNotMet
	}
	ws.capacity = capacity
	ws.events = append(ws.events, ev.StorageShrunk{
		WarehouseID:   ws.WarehouseID(),
		ShrunkStorage: capacity,
	})
	return nil
}

func (ws *WarehouseStorage) Events() []dm.Event {
	return ws.events
}

func (ws *WarehouseStorage) clearEvents() {
	ws.events = nil
}

type StorageConstructor struct{}

func (StorageConstructor) NewWarehouseStorage(
	warehouseID dm.WarehouseID,
	capacity int,
) WarehouseStorage {
	warehouseStorageID := dm.IDConstructor{}.NewWarehouseStorageID()
	queue := Queue{reservations: []Reservation{}}
	return WarehouseStorage{
		id:             warehouseStorageID,
		warehouseID:    warehouseID,
		capacity:       capacity,
		claimedStorage: 0,
		queue:          queue,
	}
}

type StorageBuilder struct{}

func (StorageBuilder) NewWarehouseStorageBuilder() WarehouseStorageBuilder {
	queue := Queue{reservations: []Reservation{}}
	warehouseStorage := WarehouseStorage{
		queue: queue,
	}
	return WarehouseStorageBuilder{warehouseStorage: &warehouseStorage}
}

type WarehouseStorageBuilder struct {
	warehouseStorage *WarehouseStorage
}

func (wsb WarehouseStorageBuilder) GetWarehouseStorage() *WarehouseStorage {
	return wsb.warehouseStorage
}

func (wsb WarehouseStorageBuilder) SetID(id string) WarehouseStorageBuilder {
	wsb.warehouseStorage.id = dm.IDConstructor{}.NewWarehouseStorageIDFromStr(id)
	return wsb
}

func (wsb WarehouseStorageBuilder) SetWarehouseID(id string) WarehouseStorageBuilder {
	wsb.warehouseStorage.warehouseID = dm.IDConstructor{}.NewWarehouseIDFromStr(id)
	return wsb
}

func (wsb WarehouseStorageBuilder) SetCapacity(capacity int) WarehouseStorageBuilder {
	wsb.warehouseStorage.capacity = capacity
	return wsb
}

func (wsb WarehouseStorageBuilder) SetClaimedStorage(claimedStorage int) WarehouseStorageBuilder {
	wsb.warehouseStorage.claimedStorage = claimedStorage
	return wsb
}

func (wsb WarehouseStorageBuilder) AddFileReservation(id string, size int) WarehouseStorageBuilder {
	fileID := dm.IDConstructor{}.NewFileIDFromStr(id)
	file := File{ID: fileID, Size: size}
	wsb.warehouseStorage.queue.Add(file)
	return wsb
}

type WarehouseStorageProxy struct {
	warehouseStorage *WarehouseStorage
}

func (wsp *WarehouseStorageProxy) GetID() string {
	return wsp.warehouseStorage.id.Str()
}

func (wsp *WarehouseStorageProxy) GetWarehouseID() string {
	return wsp.warehouseStorage.warehouseID.Str()
}

func (wsp *WarehouseStorageProxy) GetCapacity() int {
	return wsp.warehouseStorage.capacity
}

func (wsp *WarehouseStorageProxy) GetClaimedStorage() int {
	return wsp.warehouseStorage.claimedStorage
}

func (wsp *WarehouseStorageProxy) GetFileReservations() []FileReservation {
	var files []FileReservation
	for _, reservation := range wsp.warehouseStorage.queue.reservations {
		fileReservation := FileReservation{
			ID:   reservation.File.ID.Str(),
			Size: reservation.File.Size,
		}
		files = append(files, fileReservation)
	}
	return files
}

type FileReservation struct {
	ID   string
	Size int
}

func NewWarehouseStorageProxy(warehouseStorage *WarehouseStorage) WarehouseStorageProxy {
	return WarehouseStorageProxy{warehouseStorage}
}
