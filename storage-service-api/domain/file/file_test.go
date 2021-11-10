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
	sut := NewWarehouseFile(size, format)

	err := sut.RequestVerification()

	events := ut.GetFileVerificationRequestedEvents(sut.Events())
	assert.Nil(t, err)
	assert.Equal(t, st.Verifying, sut.Status())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.FileVerificationRequested{
			WarehouseID: sut.WarehouseID(),
			FileID:      sut.ID(),
		},
	)
}

func TestFileAcceptance(t *testing.T) {
	size, format := 10, CSV
	sut := NewWarehouseFile(size, format)
	sut.RequestVerification()
	sut.clearEvents()

	err := sut.Accept()

	events := ut.GetFileAcceptedEvents(sut.Events())
	assert.Nil(t, err)
	assert.Equal(t, st.Accepted, sut.Status())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.FileAccepted{
			WarehouseID: sut.WarehouseID(),
			FileID:      sut.ID(),
		},
	)
}

func TestInstantiatedFileAcceptance(t *testing.T) {
	size, format := 10, CSV
	sut := NewWarehouseFile(size, format)

	err := sut.Accept()

	events := ut.GetInstantiatedFileAcceptanceAttemptedEvents(sut.Events())
	assert.Equal(t, er.InstantiatedFileAcceptanceAttempted, err)
	assert.NotEqual(t, st.Accepted, sut.Status())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.InstantiatedFileAcceptanceAttempted{
			WarehouseID: sut.WarehouseID(),
			FileID:      sut.ID(),
		},
	)
}

func TestRemoveFile(t *testing.T) {
	size, format := 10, CSV
	sut := NewWarehouseFile(size, format)
	sut.RequestVerification()
	sut.Accept()
	sut.clearEvents()

	err := sut.Remove()

	events := ut.GetFileRemovedEvents(sut.Events())
	assert.Nil(t, err)
	assert.Equal(t, st.Removed, sut.Status())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.FileRemoved{
			WarehouseID: sut.WarehouseID(),
			FileID:      sut.ID(),
		},
	)
}

func TestUnverifiedFileRemoval(t *testing.T) {
	size, format := 10, CSV
	sut := NewWarehouseFile(size, format)
	sut.RequestVerification()
	sut.clearEvents()

	err := sut.Remove()

	events := ut.GetUnacceptedFileRemovalAttemptedEvents(sut.Events())
	assert.Equal(t, er.UnacceptedFileRemovalAttempted, err)
	assert.Equal(t, st.Verifying, sut.Status())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.UnacceptedFileRemovalAttempted{
			WarehouseID: sut.WarehouseID(),
			FileID:      sut.ID(),
		},
	)
}

func TestInstantiatedFileRemoval(t *testing.T) {
	size, format := 10, CSV
	sut := NewWarehouseFile(size, format)

	err := sut.Remove()

	events := ut.GetUnacceptedFileRemovalAttemptedEvents(sut.Events())
	assert.Equal(t, er.UnacceptedFileRemovalAttempted, err)
	assert.Equal(t, st.Instantiated, sut.Status())
	assert.Len(t, events, 1)
	assert.Contains(
		t,
		events,
		ev.UnacceptedFileRemovalAttempted{
			WarehouseID: sut.WarehouseID(),
			FileID:      sut.ID(),
		},
	)
}
