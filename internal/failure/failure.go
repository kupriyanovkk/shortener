package failure

import "errors"

// ErrConflict for data conflict cases in data base
var ErrConflict = errors.New("data conflict")

// ErrEmptyOrigURL for empty original URL case
var ErrEmptyOrigURL = errors.New("original URL cannot be empty")
