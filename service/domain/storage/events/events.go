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

type StorageExpanded struct {
	WarehouseID     dm.WarehouseID
	ExpandedStorage int
}

type MinimumStorageExpansionNotMet struct {
	WarehouseID dm.WarehouseID
	Capacity    int
}

type StorageShrunk struct {
	WarehouseID   dm.WarehouseID
	ShrunkStorage int
}

type MaximumStorageContractionNotMet struct {
	WarehouseID dm.WarehouseID
	Capacity    int
}

type MinimumStorageContractionNotMet struct {
	WarehouseID dm.WarehouseID
	Capacity    int
}

func (StorageReserved) IsEvent()                           {}
func (AvailableStorageExceeded) IsEvent()                  {}
func (StorageReservationDuplicated) IsEvent()              {}
func (ReservedStorageCommitted) IsEvent()                  {}
func (UnreservedStorageCommitted) IsEvent()                {}
func (StorageUnreserved) IsEvent()                         {}
func (NonexistentStorageReservationUnreserved) IsEvent()   {}
func (StorageFreed) IsEvent()                              {}
func (FreedStorageExceededClaimedStorage) IsEvent()        {}
func (FreeingUncommittedStorageAttempted) IsEvent()        {}
func (StorageExpanded) IsEvent()                           {}
func (MinimumStorageExpansionNotMet) IsEvent()             {}
func (StorageShrunk) IsEvent()                             {}
func (MaximumStorageContractionNotMet) IsEvent()           {}
func (MinimumStorageContractionNotMet) IsEvent()           {}
func (StorageReserved) IsMessage()                         {}
func (AvailableStorageExceeded) IsMessage()                {}
func (StorageReservationDuplicated) IsMessage()            {}
func (ReservedStorageCommitted) IsMessage()                {}
func (UnreservedStorageCommitted) IsMessage()              {}
func (StorageUnreserved) IsMessage()                       {}
func (NonexistentStorageReservationUnreserved) IsMessage() {}
func (StorageFreed) IsMessage()                            {}
func (FreedStorageExceededClaimedStorage) IsMessage()      {}
func (FreeingUncommittedStorageAttempted) IsMessage()      {}
func (StorageExpanded) IsMessage()                         {}
func (MinimumStorageExpansionNotMet) IsMessage()           {}
func (StorageShrunk) IsMessage()                           {}
func (MaximumStorageContractionNotMet) IsMessage()         {}
func (MinimumStorageContractionNotMet) IsMessage()         {}
