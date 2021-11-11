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
	warehouseID := dm.IDConstructor{}.NewWarehouseID()
	constructor := StorageConstructor{}
	warehouseStorage := constructor.NewWarehouseStorage(warehouseID, capacity)
	return warehouseStorage
}

func NewFile(size int) File {
	return File{
		ID:   dm.IDConstructor{}.NewFileID(),
		Size: size,
	}
}

func TestStorageReservation(t *testing.T) {
	capacity := 100
	sut := NewWarehouseStorage(capacity)
	file := NewFile(10)

	err := sut.Reserve(file)

	events := ut.GetStorageReservedEvents(sut.Events())
	assert.Nil(t, err)
	assert.Equal(t, 90, sut.AvailableStorage())
	assert.Equal(t, 10, sut.ReservedStorage())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.StorageReserved{
			WarehouseID: sut.WarehouseID(),
			FileID:      file.ID,
		},
	)
}

func TestExceedingStorageLimit(t *testing.T) {
	capacity := 100
	sut := NewWarehouseStorage(capacity)
	file := NewFile(10)
	otherFile := NewFile(100)
	sut.Reserve(file)
	sut.clearEvents()

	err := sut.Reserve(otherFile)

	events := ut.GetAvailableStorageExceededEvents(sut.Events())
	assert.Equal(t, er.AvailableStorageExceeded, err)
	assert.Equal(t, 90, sut.AvailableStorage())
	assert.Equal(t, 10, sut.ReservedStorage())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.AvailableStorageExceeded{
			WarehouseID:      sut.WarehouseID(),
			FileID:           otherFile.ID,
			AvailableStorage: 90,
		},
	)
}

func TestDuplicateStorageReservations(t *testing.T) {
	capacity := 100
	sut := NewWarehouseStorage(capacity)
	file := NewFile(10)
	sameFile := file
	sut.Reserve(file)
	sut.clearEvents()

	err := sut.Reserve(sameFile)

	events := ut.GetStorageReservationDuplicatedEvents(sut.Events())
	assert.Equal(t, er.StorageReservationDuplicated, err)
	assert.Equal(t, 90, sut.AvailableStorage())
	assert.Equal(t, 10, sut.ReservedStorage())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.StorageReservationDuplicated{
			WarehouseID: sut.WarehouseID(),
			FileID:      file.ID,
		},
	)
}

func TestCommittingReservation(t *testing.T) {
	capacity := 100
	sut := NewWarehouseStorage(capacity)
	file := NewFile(10)
	sut.Reserve(file)
	sut.clearEvents()

	err := sut.Commit(file.ID)

	events := ut.GetReservedStorageCommittedEvents(sut.Events())
	assert.Nil(t, err)
	assert.Equal(t, 90, sut.AvailableStorage())
	assert.Equal(t, 0, sut.ReservedStorage())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.ReservedStorageCommitted{
			WarehouseID: sut.WarehouseID(),
			FileID:      file.ID,
		},
	)
}

func TestCommittingUnreservedStorage(t *testing.T) {
	capacity := 100
	sut := NewWarehouseStorage(capacity)
	fileID := dm.IDConstructor{}.NewFileID()

	err := sut.Commit(fileID)

	events := ut.GetUnreservedStorageCommittedEvents(sut.Events())
	assert.Equal(t, er.UnreservedStorageCommitted, err)
	assert.Equal(t, 100, sut.AvailableStorage())
	assert.Equal(t, 0, sut.ReservedStorage())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.UnreservedStorageCommitted{
			WarehouseID: sut.WarehouseID(),
			FileID:      fileID,
		},
	)
}

func TestUnreservedStorage(t *testing.T) {
	capacity := 100
	sut := NewWarehouseStorage(capacity)
	file := NewFile(10)
	sut.Reserve(file)
	sut.clearEvents()

	err := sut.Unreserve(file.ID)

	events := ut.GetStorageUnreservedEvents(sut.Events())
	assert.Nil(t, err)
	assert.Equal(t, 100, sut.AvailableStorage())
	assert.Equal(t, 0, sut.ReservedStorage())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.StorageUnreserved{
			WarehouseID:       sut.WarehouseID(),
			FileID:            file.ID,
			UnreservedStorage: file.Size,
		},
	)
}

func TestUnreservingNonexistentStorageReservation(t *testing.T) {
	capacity := 100
	sut := NewWarehouseStorage(capacity)
	fileID := dm.IDConstructor{}.NewFileID()

	err := sut.Unreserve(fileID)

	events := ut.GetNonexistentStorageReservationUnreservedEvents(
		sut.Events(),
	)
	assert.Equal(t, er.NonexistentStorageReservationUnreserved, err)
	assert.Equal(t, 100, sut.AvailableStorage())
	assert.Equal(t, 0, sut.ReservedStorage())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.NonexistentStorageReservationUnreserved{
			WarehouseID: sut.WarehouseID(),
			FileID:      fileID,
		},
	)
}

func TestFreeingStorage(t *testing.T) {
	capacity := 100
	sut := NewWarehouseStorage(capacity)
	file := NewFile(10)
	sut.Reserve(file)
	sut.Commit(file.ID)
	sut.clearEvents()

	err := sut.Free(file)

	events := ut.GetStorageFreedEvents(sut.Events())
	assert.Nil(t, err)
	assert.Equal(t, 100, sut.AvailableStorage())
	assert.Equal(t, 0, sut.ReservedStorage())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.StorageFreed{
			WarehouseID:  sut.WarehouseID(),
			FileID:       file.ID,
			FreedStorage: file.Size,
		},
	)
}

func TestFreedStorageExceededClaimedStorage(t *testing.T) {
	capacity := 100
	sut := NewWarehouseStorage(capacity)
	file, largeFile := NewFile(10), NewFile(15)
	sut.Reserve(file)
	sut.Commit(file.ID)
	sut.clearEvents()

	err := sut.Free(largeFile)

	events := ut.GetFreedStorageExceededClaimedStorageEvents(sut.Events())
	assert.Equal(t, er.FreedStorageExceededClaimedStorage, err)
	assert.Equal(t, 90, sut.AvailableStorage())
	assert.Equal(t, 0, sut.ReservedStorage())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.FreedStorageExceededClaimedStorage{
			WarehouseID:    sut.WarehouseID(),
			FileID:         largeFile.ID,
			ClaimedStorage: 10,
		},
	)
}

func TestFreeingUncommittedStorage(t *testing.T) {
	capacity := 100
	sut := NewWarehouseStorage(capacity)
	file := NewFile(10)
	sut.Reserve(file)
	sut.clearEvents()

	err := sut.Free(file)

	events := ut.GetFreeingUncommittedStorageAttemptedEvents(
		sut.Events(),
	)
	assert.Equal(t, er.FreeingUncommittedStorageAttempted, err)
	assert.Equal(t, 90, sut.AvailableStorage())
	assert.Equal(t, 10, sut.ReservedStorage())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.FreeingUncommittedStorageAttempted{
			WarehouseID: sut.WarehouseID(),
			FileID:      file.ID,
		},
	)
}

func TestExpandStorage(t *testing.T) {
	capacity, expandedCapacity := 100, 1000
	sut := NewWarehouseStorage(capacity)
	file, otherFile := NewFile(10), NewFile(15)
	sut.Reserve(file)
	sut.Reserve(otherFile)
	sut.Commit(otherFile.ID)
	sut.clearEvents()

	err := sut.Expand(expandedCapacity)

	events := ut.GetStorageExpandedEvents(sut.Events())
	assert.Nil(t, err)
	assert.Equal(t, 975, sut.AvailableStorage())
	assert.Equal(t, 10, sut.ReservedStorage())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.StorageExpanded{
			WarehouseID:     sut.WarehouseID(),
			ExpandedStorage: 1000,
		},
	)
}

func TestMinimalStorageExpansion(t *testing.T) {
	capacity, smallerCapacity := 100, 50
	sut := NewWarehouseStorage(capacity)
	file, otherFile := NewFile(10), NewFile(15)
	sut.Reserve(file)
	sut.Reserve(otherFile)
	sut.Commit(otherFile.ID)
	sut.clearEvents()

	err := sut.Expand(smallerCapacity)

	events := ut.GetMinimumStorageExpansionNotMetEvents(sut.Events())
	assert.Equal(t, er.MinimumStorageExpansionNotMet, err)
	assert.Equal(t, 75, sut.AvailableStorage())
	assert.Equal(t, 10, sut.ReservedStorage())
	assert.Equal(t, 100, sut.Capacity())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.MinimumStorageExpansionNotMet{
			WarehouseID: sut.WarehouseID(),
			Capacity:    100,
		},
	)
}

func TestShrinkStorage(t *testing.T) {
	capacity, shrunkCapacity := 100, 50
	sut := NewWarehouseStorage(capacity)
	file, otherFile := NewFile(10), NewFile(15)
	sut.Reserve(file)
	sut.Reserve(otherFile)
	sut.Commit(otherFile.ID)
	sut.clearEvents()

	err := sut.Shrink(shrunkCapacity)

	events := ut.GetStorageShrunkEvents(sut.Events())
	assert.Nil(t, err)
	assert.Equal(t, 25, sut.AvailableStorage())
	assert.Equal(t, 10, sut.ReservedStorage())
	assert.Equal(t, 50, sut.Capacity())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.StorageShrunk{
			WarehouseID:   sut.WarehouseID(),
			ShrunkStorage: 50,
		},
	)
}

func TestMaximumStorageContraction(t *testing.T) {
	capacity, largerCapacity := 100, 1000
	sut := NewWarehouseStorage(capacity)
	file, otherFile := NewFile(10), NewFile(15)
	sut.Reserve(file)
	sut.Reserve(otherFile)
	sut.Commit(otherFile.ID)
	sut.clearEvents()

	err := sut.Shrink(largerCapacity)

	events := ut.GetMaximumStorageContractionNotMetEvents(
		sut.Events(),
	)
	assert.Equal(t, er.MaximumStorageContractionNotMet, err)
	assert.Equal(t, 75, sut.AvailableStorage())
	assert.Equal(t, 10, sut.ReservedStorage())
	assert.Equal(t, 100, sut.Capacity())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.MaximumStorageContractionNotMet{
			WarehouseID: sut.WarehouseID(),
			Capacity:    100,
		},
	)
}

func TestMinimumStorageContraction(t *testing.T) {
	capacity, invalidCapacity := 100, 20
	sut := NewWarehouseStorage(capacity)
	file, otherFile := NewFile(10), NewFile(15)
	sut.Reserve(file)
	sut.Reserve(otherFile)
	sut.Commit(otherFile.ID)
	sut.clearEvents()

	err := sut.Shrink(invalidCapacity)

	events := ut.GetMinimumStorageContractionNotMetEvents(
		sut.Events(),
	)
	assert.Equal(t, er.MinimumStorageContractionNotMet, err)
	assert.Equal(t, 75, sut.AvailableStorage())
	assert.Equal(t, 10, sut.ReservedStorage())
	assert.Equal(t, 100, sut.Capacity())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.MinimumStorageContractionNotMet{
			WarehouseID: sut.WarehouseID(),
			Capacity:    100,
		},
	)
}
