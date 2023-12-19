package model

type PermissionType string

const (
	DefaultUser PermissionType = "default"
	AdminUser   PermissionType = "admin"
)
