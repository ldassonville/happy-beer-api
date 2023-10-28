package client

import (
	"errors"
	"fmt"
	"net/http"
)

type StatusReason string

const (
	StatusReasonUnknown = ""

	StatusReasonNotFound = "NotFound"

	StatusReasonUnauthorized = "Unauthorized"
	StatusReasonForbidden    = "Forbidden"
	StatusInternalError      = "InternalError"
	StatusBadRequest         = "BadRequest"
	StatusServiceUnavailable = "ServiceUnavailable"
	StatusConflict           = "Conflict"
	StatusPreconditionFailed = "PreconditionFailed"
)

var knownReasons = map[StatusReason]struct{}{
	StatusReasonNotFound:     {},
	StatusReasonUnknown:      {},
	StatusReasonForbidden:    {},
	StatusReasonUnauthorized: {},
	StatusInternalError:      {},
	StatusServiceUnavailable: {},
	StatusPreconditionFailed: {},
	StatusConflict:           {},
}

type APIError struct {
	// Code is the HTTP status code
	Code int
	// Reason is the machine-readable reason for the error
	Reason StatusReason
	// Message is the human-readable description of the error
	Message string
}

var err error = &APIError{}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Reason, e.Message)
}

func IsNotFound(err error) bool {
	reason, code := reasonAndCodeForError(err)
	if reason == StatusReasonNotFound {
		return true
	}
	if _, ok := knownReasons[reason]; !ok && code == http.StatusNotFound {
		return true
	}
	return false
}

func reasonAndCodeForError(err error) (StatusReason, int) {
	var apiError *APIError
	if errors.As(err, &apiError) {
		return apiError.Reason, apiError.Code
	}
	return StatusReasonUnknown, 0
}

func ReasonForError(err error) StatusReason {

	reason, _ := reasonAndCodeForError(err)
	return reason
}
