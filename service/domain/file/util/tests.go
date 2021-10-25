package util

import (
	dm "storage.g128.rvlz.io/domain"
	ev "storage.g128.rvlz.io/domain/file/events"
)

func GetFileVerificationRequestedEvents(events []dm.Event) []ev.FileVerificationRequested {
	targetEvents := []ev.FileVerificationRequested{}
	for _, event := range events {
		targetEvent, ok := event.(ev.FileVerificationRequested)
		if ok {
			targetEvents = append(targetEvents, targetEvent)
		}
	}
	return targetEvents
}
