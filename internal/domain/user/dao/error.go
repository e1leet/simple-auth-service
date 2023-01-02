package dao

import "errors"

var ErrUsernameAlreadyUsed = errors.New("username already used")
var ErrUserNotFound = errors.New("user not found")
