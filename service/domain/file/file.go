package file

import (
	dm "storage.g128.rvlz.io/domain"
	ev "storage.g128.rvlz.io/domain/file/events"
	st "storage.g128.rvlz.io/domain/file/status"
)

type Format struct {
	format string
}

var (
	CSV = Format{"CSV"}
)

type WarehouseFile struct {
	id          dm.FileID
	warehouseID dm.WarehouseID
	size        int
	format      Format
	status      st.Status
	events      []dm.Event
}

func (wf *WarehouseFile) ID() dm.FileID {
	return wf.id
}

func (wf *WarehouseFile) WarehouseID() dm.WarehouseID {
	return wf.warehouseID
}

func (wf *WarehouseFile) Status() st.Status {
	return wf.status
}

func (wf *WarehouseFile) RequestVerification() error {
	wf.status = st.Verifying
	wf.events = append(wf.events, ev.FileVerificationRequested{
		WarehouseID: wf.WarehouseID(),
		FileID:      wf.ID(),
	})
	return nil
}

func (wf *WarehouseFile) Events() []dm.Event {
	return wf.events
}

type FileFactory struct{}

func (factory FileFactory) NewWarehouseFile(
	warehouseID dm.WarehouseID,
	size int,
	format Format,
) WarehouseFile {
	fileID := dm.IDFactory{}.NewFileID()
	return WarehouseFile{
		id:          fileID,
		warehouseID: warehouseID,
		size:        size,
		format:      format,
		status:      st.Instantiated,
	}
}
