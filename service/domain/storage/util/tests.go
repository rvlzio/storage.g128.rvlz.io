package util

import (
	dm "storage.g128.rvlz.io/domain"
	ev "storage.g128.rvlz.io/domain/storage/events"
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

func GetUnreservedStorageCommittedEvents(events []dm.Event) []ev.UnreservedStorageCommitted {
	targetEvents := []ev.UnreservedStorageCommitted{}
	for _, event := range events {
		targetEvent, ok := event.(ev.UnreservedStorageCommitted)
		if ok {
			targetEvents = append(targetEvents, targetEvent)
		}
	}
	return targetEvents
}

func GetStorageUnreservedEvents(events []dm.Event) []ev.StorageUnreserved {
	targetEvents := []ev.StorageUnreserved{}
	for _, event := range events {
		targetEvent, ok := event.(ev.StorageUnreserved)
		if ok {
			targetEvents = append(targetEvents, targetEvent)
		}
	}
	return targetEvents
}

func GetNonexistentStorageReservationUnreservedEvents(
	events []dm.Event,
) []ev.NonexistentStorageReservationUnreserved {
	targetEvents := []ev.NonexistentStorageReservationUnreserved{}
	for _, event := range events {
		targetEvent, ok := event.(ev.NonexistentStorageReservationUnreserved)
		if ok {
			targetEvents = append(targetEvents, targetEvent)
		}
	}
	return targetEvents
}