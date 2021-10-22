package storage

import (
	"github.com/stretchr/testify/assert"
	dm "storage.g128.rvlz.io/domain"
	er "storage.g128.rvlz.io/domain/storage/errors"
	ev "storage.g128.rvlz.io/domain/storage/events"
	ut "storage.g128.rvlz.io/domain/storage/util"
	"testing"
)

func TestStorageReservation(t *testing.T) {
	warehouseID, capacity := dm.IDFactory{}.NewWarehouseID(), 100
	factory := StorageFactory{}
	warehouseStorage := factory.NewWarehouseStorage(warehouseID, capacity)
	file := File{
		ID:   dm.IDFactory{}.NewFileID(),
		Size: 10,
	}

	err := warehouseStorage.Reserve(file)

	events := ut.GetStorageReservedEvents(warehouseStorage.Events())
	assert.Nil(t, err)
	assert.Equal(t, 90, warehouseStorage.AvailableStorage())
	assert.Equal(t, 10, warehouseStorage.ReservedStorage())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.StorageReserved{
			WarehouseID: warehouseID,
			FileID:      file.ID,
		},
	)
}

func TestExceedingStorageLimit(t *testing.T) {
	warehouseID, capacity := dm.IDFactory{}.NewWarehouseID(), 100
	factory := StorageFactory{}
	warehouseStorage := factory.NewWarehouseStorage(warehouseID, capacity)
	file := File{
		ID:   dm.IDFactory{}.NewFileID(),
		Size: 10,
	}
	otherFile := File{
		ID:   dm.IDFactory{}.NewFileID(),
		Size: 100,
	}
	warehouseStorage.Reserve(file)
	warehouseStorage.clearEvents()

	err := warehouseStorage.Reserve(otherFile)

	events := ut.GetAvailableStorageExceededEvents(warehouseStorage.Events())
	assert.Equal(t, er.AvailableStorageExceeded, err)
	assert.Equal(t, 90, warehouseStorage.AvailableStorage())
	assert.Equal(t, 10, warehouseStorage.ReservedStorage())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.AvailableStorageExceeded{
			WarehouseID:      warehouseID,
			FileID:           otherFile.ID,
			AvailableStorage: 90,
		},
	)
}

func TestDuplicateStorageReservations(t *testing.T) {
	warehouseID, capacity := dm.IDFactory{}.NewWarehouseID(), 100
	factory := StorageFactory{}
	warehouseStorage := factory.NewWarehouseStorage(warehouseID, capacity)
	file := File{
		ID:   dm.IDFactory{}.NewFileID(),
		Size: 10,
	}
	sameFile := file
	warehouseStorage.Reserve(file)
	warehouseStorage.clearEvents()

	err := warehouseStorage.Reserve(sameFile)

	events := ut.GetStorageReservationDuplicatedEvents(warehouseStorage.Events())
	assert.Equal(t, er.StorageReservationDuplicated, err)
	assert.Equal(t, 90, warehouseStorage.AvailableStorage())
	assert.Equal(t, 10, warehouseStorage.ReservedStorage())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.StorageReservationDuplicated{
			WarehouseID: warehouseID,
			FileID:      file.ID,
		},
	)
}

func TestCommittingReservation(t *testing.T) {
	warehouseID, capacity := dm.IDFactory{}.NewWarehouseID(), 100
	factory := StorageFactory{}
	warehouseStorage := factory.NewWarehouseStorage(warehouseID, capacity)
	file := File{
		ID:   dm.IDFactory{}.NewFileID(),
		Size: 10,
	}
	warehouseStorage.Reserve(file)
	warehouseStorage.clearEvents()

	err := warehouseStorage.Commit(file.ID)

	events := ut.GetReservedStorageCommittedEvents(warehouseStorage.Events())
	assert.Nil(t, err)
	assert.Equal(t, 90, warehouseStorage.AvailableStorage())
	assert.Equal(t, 0, warehouseStorage.ReservedStorage())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.ReservedStorageCommitted{
			WarehouseID: warehouseID,
			FileID:      file.ID,
		},
	)
}

func TestCommittingUnreservedStorage(t *testing.T) {
	warehouseID, capacity := dm.IDFactory{}.NewWarehouseID(), 100
	factory := StorageFactory{}
	warehouseStorage := factory.NewWarehouseStorage(warehouseID, capacity)
	fileID := dm.IDFactory{}.NewFileID()

	err := warehouseStorage.Commit(fileID)

	events := ut.GetUnreservedStorageCommittedEvents(warehouseStorage.Events())
	assert.Equal(t, er.UnreservedStorageCommitted, err)
	assert.Equal(t, 100, warehouseStorage.AvailableStorage())
	assert.Equal(t, 0, warehouseStorage.ReservedStorage())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.UnreservedStorageCommitted{
			WarehouseID: warehouseID,
			FileID:      fileID,
		},
	)
}
