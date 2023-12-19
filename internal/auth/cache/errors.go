package cache

import "errors"

var (
	ErrAccountAlreadyExists    = errors.New("the account already exists")
	ErrAccountNotCreated       = errors.New("the account not created")
	ErrPermissionAlreadyExists = errors.New("the permission already exists")
	ErrPermissionNotFound      = errors.New("permissions not found")
)
