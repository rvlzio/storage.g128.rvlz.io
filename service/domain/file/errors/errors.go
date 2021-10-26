package errors

import (
	"errors"
)

var (
	FileAcceptedBeforeVerificationRequest = errors.New(
		"file_accepted_before_verification_request",
	)
	UnacceptedFileRemovalAttempted = errors.New(
		"unaccepted_file_removal_attempted",
	)
)
