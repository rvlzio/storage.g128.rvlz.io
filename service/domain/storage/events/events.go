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

func (StorageReserved) IsEvent()              {}
func (AvailableStorageExceeded) IsEvent()     {}
func (StorageReservationDuplicated) IsEvent() {}
func (ReservedStorageCommitted) IsEvent()     {}
func (UnreservedStorageCommitted) IsEvent()   {}
