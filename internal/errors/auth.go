package errors

import "errors"

var ErrInvalidCredentials = errors.New("invalid login or password")
var ErrUserIsAlreadyExists = errors.New("login already taken")
