package storage

import (
	"github.com/stretchr/testify/assert"
	dm "storage.g128.rvlz.io/domain"
	er "storage.g128.rvlz.io/domain/storage/errors"
	ev "storage.g128.rvlz.io/domain/storage/events"
	ut "storage.g128.rvlz.io/domain/storage/util"
	"testing"
)

func NewWarehouseStorage(capacity int) WarehouseStorage {
	warehouseID := dm.IDFactory{}.NewWarehouseID()
	factory := StorageFactory{}
	warehouseStorage := factory.NewWarehouseStorage(warehouseID, capacity)
	return warehouseStorage
}

func NewFile(size int) File {
	return File{
		ID:   dm.IDFactory{}.NewFileID(),
		Size: size,
	}
}

func TestStorageReservation(t *testing.T) {
	capacity := 100
	warehouseStorage := NewWarehouseStorage(capacity)
	file := NewFile(10)

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
			WarehouseID: warehouseStorage.WarehouseID(),
			FileID:      file.ID,
		},
	)
}

func TestExceedingStorageLimit(t *testing.T) {
	capacity := 100
	warehouseStorage := NewWarehouseStorage(capacity)
	file := NewFile(10)
	otherFile := NewFile(100)
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
			WarehouseID:      warehouseStorage.WarehouseID(),
			FileID:           otherFile.ID,
			AvailableStorage: 90,
		},
	)
}

func TestDuplicateStorageReservations(t *testing.T) {
	capacity := 100
	warehouseStorage := NewWarehouseStorage(capacity)
	file := NewFile(10)
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
			WarehouseID: warehouseStorage.WarehouseID(),
			FileID:      file.ID,
		},
	)
}

func TestCommittingReservation(t *testing.T) {
	capacity := 100
	warehouseStorage := NewWarehouseStorage(capacity)
	file := NewFile(10)
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
			WarehouseID: warehouseStorage.WarehouseID(),
			FileID:      file.ID,
		},
	)
}

func TestCommittingUnreservedStorage(t *testing.T) {
	capacity := 100
	warehouseStorage := NewWarehouseStorage(capacity)
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
			WarehouseID: warehouseStorage.WarehouseID(),
			FileID:      fileID,
		},
	)
}

func TestUnreservedStorage(t *testing.T) {
	capacity := 100
	warehouseStorage := NewWarehouseStorage(capacity)
	file := NewFile(10)
	warehouseStorage.Reserve(file)
	warehouseStorage.clearEvents()

	err := warehouseStorage.Unreserve(file.ID)

	events := ut.GetStorageUnreservedEvents(warehouseStorage.Events())
	assert.Nil(t, err)
	assert.Equal(t, 100, warehouseStorage.AvailableStorage())
	assert.Equal(t, 0, warehouseStorage.ReservedStorage())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.StorageUnreserved{
			WarehouseID:       warehouseStorage.WarehouseID(),
			FileID:            file.ID,
			UnreservedStorage: file.Size,
		},
	)
}

func TestUnreservingNonexistentStorageReservation(t *testing.T) {
	capacity := 100
	warehouseStorage := NewWarehouseStorage(capacity)
	fileID := dm.IDFactory{}.NewFileID()

	err := warehouseStorage.Unreserve(fileID)

	events := ut.GetNonexistentStorageReservationUnreservedEvents(
		warehouseStorage.Events(),
	)
	assert.Equal(t, er.NonexistentStorageReservationUnreserved, err)
	assert.Equal(t, 100, warehouseStorage.AvailableStorage())
	assert.Equal(t, 0, warehouseStorage.ReservedStorage())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.NonexistentStorageReservationUnreserved{
			WarehouseID: warehouseStorage.WarehouseID(),
			FileID:      fileID,
		},
	)
}

func TestFreeingStorage(t *testing.T) {
	capacity := 100
	warehouseStorage := NewWarehouseStorage(capacity)
	file := NewFile(10)
	warehouseStorage.Reserve(file)
	warehouseStorage.Commit(file.ID)
	warehouseStorage.clearEvents()

	err := warehouseStorage.Free(file)

	events := ut.GetStorageFreedEvents(warehouseStorage.Events())
	assert.Nil(t, err)
	assert.Equal(t, 100, warehouseStorage.AvailableStorage())
	assert.Equal(t, 0, warehouseStorage.ReservedStorage())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.StorageFreed{
			WarehouseID:  warehouseStorage.WarehouseID(),
			FileID:       file.ID,
			FreedStorage: file.Size,
		},
	)
}

func TestFreedStorageExceedAvailability(t *testing.T) {
	capacity := 100
	warehouseStorage := NewWarehouseStorage(capacity)
	file, largeFile := NewFile(10), NewFile(95)
	warehouseStorage.Reserve(file)
	warehouseStorage.clearEvents()

	err := warehouseStorage.Free(largeFile)

	events := ut.GetFreedStorageExceededAvailabilityEvents(warehouseStorage.Events())
	assert.Equal(t, er.FreedStorageExceededAvailability, err)
	assert.Equal(t, 90, warehouseStorage.AvailableStorage())
	assert.Equal(t, 10, warehouseStorage.ReservedStorage())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.FreedStorageExceededAvailability{
			WarehouseID:      warehouseStorage.WarehouseID(),
			FileID:           largeFile.ID,
			AvailableStorage: 90,
		},
	)
}

func TestFreeingUncommittedStorage(t *testing.T) {
	capacity := 100
	warehouseStorage := NewWarehouseStorage(capacity)
	file := NewFile(10)
	warehouseStorage.Reserve(file)
	warehouseStorage.clearEvents()

	err := warehouseStorage.Free(file)

	events := ut.GetFreeingUncommittedStorageAttemptedEvents(
		warehouseStorage.Events(),
	)
	assert.Equal(t, er.FreeingUncommittedStorageAttempted, err)
	assert.Equal(t, 90, warehouseStorage.AvailableStorage())
	assert.Equal(t, 10, warehouseStorage.ReservedStorage())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.FreeingUncommittedStorageAttempted{
			WarehouseID: warehouseStorage.WarehouseID(),
			FileID:      file.ID,
		},
	)
}
