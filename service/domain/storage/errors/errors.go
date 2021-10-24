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
	FreedStorageExceededClaimedStorage = errors.New(
		"freed_storage_exceeded_claimed_storage",
	)
	FreeingUncommittedStorageAttempted = errors.New(
		"freeing_uncommitted_storage_attempted",
	)
)
