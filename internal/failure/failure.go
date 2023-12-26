package failure

import "errors"

var ErrConflict = errors.New("data conflict")
var ErrEmptyOrigURL = errors.New("original URL cannot be empty")
