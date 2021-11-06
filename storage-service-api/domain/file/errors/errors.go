package errors

import (
	"errors"
)

var (
	InstantiatedFileAcceptanceAttempted = errors.New(
		"instantiated_file_acceptance_attempted",
	)
	UnacceptedFileRemovalAttempted = errors.New(
		"unaccepted_file_removal_attempted",
	)
)
