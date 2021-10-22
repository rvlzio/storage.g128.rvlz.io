package storage

import (
	"github.com/stretchr/testify/assert"
	dm "storage.g128.rvlz.io/domain"
	er "storage.g128.rvlz.io/domain/storage/errors"
	ev "storage.g128.rvlz.io/domain/storage/events"
	"testing"
)

func GetStorageReservedEvents(events []dm.Event) []ev.StorageReserved {
	targetEvents := []ev.StorageReserved{}
	for _, event := range events {
		targetEvent, ok := event.(ev.StorageReserved)
		if ok {
			targetEvents = append(targetEvents, targetEvent)
		}
	}
	return targetEvents
}

func GetAvailableStorageExceededEvents(events []dm.Event) []ev.AvailableStorageExceeded {
	targetEvents := []ev.AvailableStorageExceeded{}
	for _, event := range events {
		targetEvent, ok := event.(ev.AvailableStorageExceeded)
		if ok {
			targetEvents = append(targetEvents, targetEvent)
		}
	}
	return targetEvents
}

func GetStorageReservationDuplicatedEvents(events []dm.Event) []ev.StorageReservationDuplicated {
	targetEvents := []ev.StorageReservationDuplicated{}
	for _, event := range events {
		targetEvent, ok := event.(ev.StorageReservationDuplicated)
		if ok {
			targetEvents = append(targetEvents, targetEvent)
		}
	}
	return targetEvents
}

func GetReservedStorageCommittedEvents(events []dm.Event) []ev.ReservedStorageCommitted {
	targetEvents := []ev.ReservedStorageCommitted{}
	for _, event := range events {
		targetEvent, ok := event.(ev.ReservedStorageCommitted)
		if ok {
			targetEvents = append(targetEvents, targetEvent)
		}
	}
	return targetEvents
}

func TestStorageReservation(t *testing.T) {
	warehouseID, capacity := dm.IDFactory{}.NewWarehouseID(), 100
	factory := StorageFactory{}
	warehouseStorage := factory.NewWarehouseStorage(warehouseID, capacity)
	file := File{
		ID:   dm.IDFactory{}.NewFileID(),
		Size: 10,
	}

	err := warehouseStorage.Reserve(file)

	events := GetStorageReservedEvents(warehouseStorage.Events())
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

	events := GetAvailableStorageExceededEvents(warehouseStorage.Events())
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

	events := GetStorageReservationDuplicatedEvents(warehouseStorage.Events())
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

	events := GetReservedStorageCommittedEvents(warehouseStorage.Events())
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
