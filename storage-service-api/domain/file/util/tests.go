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

func GetFileAcceptedEvents(events []dm.Event) []ev.FileAccepted {
	targetEvents := []ev.FileAccepted{}
	for _, event := range events {
		targetEvent, ok := event.(ev.FileAccepted)
		if ok {
			targetEvents = append(targetEvents, targetEvent)
		}
	}
	return targetEvents
}

func GetInstantiatedFileAcceptanceAttemptedEvents(
	events []dm.Event,
) []ev.InstantiatedFileAcceptanceAttempted {
	targetEvents := []ev.InstantiatedFileAcceptanceAttempted{}
	for _, event := range events {
		targetEvent, ok := event.(ev.InstantiatedFileAcceptanceAttempted)
		if ok {
			targetEvents = append(targetEvents, targetEvent)
		}
	}
	return targetEvents
}

func GetFileRemovedEvents(events []dm.Event) []ev.FileRemoved {
	targetEvents := []ev.FileRemoved{}
	for _, event := range events {
		targetEvent, ok := event.(ev.FileRemoved)
		if ok {
			targetEvents = append(targetEvents, targetEvent)
		}
	}
	return targetEvents
}

func GetUnacceptedFileRemovalAttemptedEvents(
	events []dm.Event,
) []ev.UnacceptedFileRemovalAttempted {
	targetEvents := []ev.UnacceptedFileRemovalAttempted{}
	for _, event := range events {
		targetEvent, ok := event.(ev.UnacceptedFileRemovalAttempted)
		if ok {
			targetEvents = append(targetEvents, targetEvent)
		}
	}
	return targetEvents
}
