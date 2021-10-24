package events

import (
	dm "storage.g128.rvlz.io/domain"
)

type StorageReserved struct {
	WarehouseID dm.WarehouseID
	FileID      dm.FileID
}

type AvailableStorageExceeded struct {
	WarehouseID      dm.WarehouseID
	FileID           dm.FileID
	AvailableStorage int
}

type StorageReservationDuplicated struct {
	WarehouseID dm.WarehouseID
	FileID      dm.FileID
}

type ReservedStorageCommitted struct {
	WarehouseID dm.WarehouseID
	FileID      dm.FileID
}

type UnreservedStorageCommitted struct {
	WarehouseID dm.WarehouseID
	FileID      dm.FileID
}

type StorageUnreserved struct {
	WarehouseID       dm.WarehouseID
	FileID            dm.FileID
	UnreservedStorage int
}

type NonexistentStorageReservationUnreserved struct {
	WarehouseID dm.WarehouseID
	FileID      dm.FileID
}

type StorageFreed struct {
	WarehouseID  dm.WarehouseID
	FileID       dm.FileID
	FreedStorage int
}

type FreedStorageExceededClaimedStorage struct {
	WarehouseID    dm.WarehouseID
	FileID         dm.FileID
	ClaimedStorage int
}

type FreeingUncommittedStorageAttempted struct {
	WarehouseID dm.WarehouseID
	FileID      dm.FileID
}

func (StorageReserved) IsEvent()                         {}
func (AvailableStorageExceeded) IsEvent()                {}
func (StorageReservationDuplicated) IsEvent()            {}
func (ReservedStorageCommitted) IsEvent()                {}
func (UnreservedStorageCommitted) IsEvent()              {}
func (StorageUnreserved) IsEvent()                       {}
func (NonexistentStorageReservationUnreserved) IsEvent() {}
func (StorageFreed) IsEvent()                            {}
func (FreedStorageExceededClaimedStorage) IsEvent()      {}
func (FreeingUncommittedStorageAttempted) IsEvent()      {}
