package errors

import "errors"

var InvalidCredentials = errors.New("Invalid login or password")
var UserIsAlreadyExists = errors.New("Login already taken")
