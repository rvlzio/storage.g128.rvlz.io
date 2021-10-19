package errors

import (
	"errors"
)

var (
	AvailableStorageExceeded = errors.New(
		"available_storage_exceeded",
	)
	StorageReservationDuplicated = errors.New(
		"storage_reservation_duplicated",
	)
)
