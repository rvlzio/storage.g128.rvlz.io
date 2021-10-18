package storage

import (
	"github.com/stretchr/testify/assert"
	dm "storage.g128.rvlz.io/domain"
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
