package errorsutil

import "errors"

var ErrDefault = errors.New("error")
var ErrEOF = errors.New("EOF")
var ErrNotFound = errors.New("not found")
var ErrNoMongoDoc = errors.New("mongo: no documents in result")
