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
	UnreservedStorageCommitted = errors.New(
		"unreserved_storage_committed",
	)
	NonexistentStorageReservationUnreserved = errors.New(
		"nonexistent_storage_reservation_unreserved",
	)
	FreedStorageExceededAvailability = errors.New(
		"freed_storage_exceeded_availability",
	)
)
