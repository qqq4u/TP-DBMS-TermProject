package models

import "errors"

var (
	ErrorConflict = errors.New("Conflict")
	ErrorNotFound = errors.New("NotFound")
	ErrorInternal = errors.New("InternalError")
)
