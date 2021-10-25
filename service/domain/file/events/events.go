package events

import (
	dm "storage.g128.rvlz.io/domain"
)

type FileVerificationRequested struct {
	WarehouseID dm.WarehouseID
	FileID      dm.FileID
}

type FileAccepted struct {
	WarehouseID dm.WarehouseID
	FileID      dm.FileID
}

func (FileVerificationRequested) IsEvent() {}
func (FileAccepted) IsEvent()              {}
