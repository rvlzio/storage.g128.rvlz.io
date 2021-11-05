package file

import (
	"github.com/stretchr/testify/assert"
	dm "storage.g128.rvlz.io/domain"
	er "storage.g128.rvlz.io/domain/file/errors"
	ev "storage.g128.rvlz.io/domain/file/events"
	st "storage.g128.rvlz.io/domain/file/status"
	ut "storage.g128.rvlz.io/domain/file/util"
	"testing"
)

func NewWarehouseFile(size int, format Format) WarehouseFile {
	warehouseID := dm.IDConstructor{}.NewWarehouseID()
	constructor := FileConstructor{}
	return constructor.NewWarehouseFile(warehouseID, size, format)
}

func TestVerificationRequest(t *testing.T) {
	size, format := 10, CSV
	warehouseFile := NewWarehouseFile(size, format)

	err := warehouseFile.RequestVerification()

	events := ut.GetFileVerificationRequestedEvents(warehouseFile.Events())
	assert.Nil(t, err)
	assert.Equal(t, st.Verifying, warehouseFile.Status())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.FileVerificationRequested{
			WarehouseID: warehouseFile.WarehouseID(),
			FileID:      warehouseFile.ID(),
		},
	)
}

func TestFileAcceptance(t *testing.T) {
	size, format := 10, CSV
	warehouseFile := NewWarehouseFile(size, format)
	warehouseFile.RequestVerification()
	warehouseFile.clearEvents()

	err := warehouseFile.Accept()

	events := ut.GetFileAcceptedEvents(warehouseFile.Events())
	assert.Nil(t, err)
	assert.Equal(t, st.Accepted, warehouseFile.Status())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.FileAccepted{
			WarehouseID: warehouseFile.WarehouseID(),
			FileID:      warehouseFile.ID(),
		},
	)
}

func TestInstantiatedFileAcceptance(t *testing.T) {
	size, format := 10, CSV
	warehouseFile := NewWarehouseFile(size, format)

	err := warehouseFile.Accept()

	events := ut.GetInstantiatedFileAcceptanceAttemptedEvents(warehouseFile.Events())
	assert.Equal(t, er.InstantiatedFileAcceptanceAttempted, err)
	assert.NotEqual(t, st.Accepted, warehouseFile.Status())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.InstantiatedFileAcceptanceAttempted{
			WarehouseID: warehouseFile.WarehouseID(),
			FileID:      warehouseFile.ID(),
		},
	)
}

func TestRemoveFile(t *testing.T) {
	size, format := 10, CSV
	warehouseFile := NewWarehouseFile(size, format)
	warehouseFile.RequestVerification()
	warehouseFile.Accept()
	warehouseFile.clearEvents()

	err := warehouseFile.Remove()

	events := ut.GetFileRemovedEvents(warehouseFile.Events())
	assert.Nil(t, err)
	assert.Equal(t, st.Removed, warehouseFile.Status())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.FileRemoved{
			WarehouseID: warehouseFile.WarehouseID(),
			FileID:      warehouseFile.ID(),
		},
	)
}

func TestUnverifiedFileRemoval(t *testing.T) {
	size, format := 10, CSV
	warehouseFile := NewWarehouseFile(size, format)
	warehouseFile.RequestVerification()
	warehouseFile.clearEvents()

	err := warehouseFile.Remove()

	events := ut.GetUnacceptedFileRemovalAttemptedEvents(warehouseFile.Events())
	assert.Equal(t, er.UnacceptedFileRemovalAttempted, err)
	assert.Equal(t, st.Verifying, warehouseFile.Status())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.UnacceptedFileRemovalAttempted{
			WarehouseID: warehouseFile.WarehouseID(),
			FileID:      warehouseFile.ID(),
		},
	)
}

func TestInstantiatedFileRemoval(t *testing.T) {
	size, format := 10, CSV
	warehouseFile := NewWarehouseFile(size, format)

	err := warehouseFile.Remove()

	events := ut.GetUnacceptedFileRemovalAttemptedEvents(warehouseFile.Events())
	assert.Equal(t, er.UnacceptedFileRemovalAttempted, err)
	assert.Equal(t, st.Instantiated, warehouseFile.Status())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.UnacceptedFileRemovalAttempted{
			WarehouseID: warehouseFile.WarehouseID(),
			FileID:      warehouseFile.ID(),
		},
	)
}
