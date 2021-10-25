package file

import (
	"github.com/stretchr/testify/assert"
	dm "storage.g128.rvlz.io/domain"
	ev "storage.g128.rvlz.io/domain/file/events"
	st "storage.g128.rvlz.io/domain/file/status"
	ut "storage.g128.rvlz.io/domain/file/util"
	"testing"
)

func NewWarehouseFile(size int, format Format) WarehouseFile {
	warehouseID := dm.IDFactory{}.NewWarehouseID()
	factory := FileFactory{}
	return factory.NewWarehouseFile(warehouseID, size, format)
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
