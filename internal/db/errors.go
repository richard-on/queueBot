package db

import "errors"

var ErrLackUserInfo = errors.New("user not initialised")

var ErrNoUserInfo = errors.New("user not found")
