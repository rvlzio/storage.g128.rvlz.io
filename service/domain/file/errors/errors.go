package errors

import (
	"errors"
)

var (
	FileAcceptedBeforeVerificationRequest = errors.New(
		"file_accepted_before_verification_request",
	)
)
