package cache

import (
	"sync"
	"vsu/internal/auth/model"
)

type AccountCache struct {
	cache map[string]string
	mx    *sync.Mutex
}

type PermissionsCache struct {
	cache map[string][]model.PermissionType
	mx    *sync.Mutex
}

func NewAccountCache() *AccountCache {
	return &AccountCache{
		cache: make(map[string]string),
		mx:    &sync.Mutex{},
	}
}

func NewPermissionsCache() *PermissionsCache {
	return &PermissionsCache{
		cache: make(map[string][]model.PermissionType),
		mx:    &sync.Mutex{},
	}
}

func (c *AccountCache) Set(login, password string) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	if _, ok := c.cache[login]; ok {
		return ErrAccountAlreadyExists
	}

	c.cache[login] = password

	return nil
}

func (c *AccountCache) GetPass(login string) (string, error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	password, ok := c.cache[login]
	if !ok {
		return "", ErrAccountNotCreated
	}

	return password, nil
}

func (c *PermissionsCache) Set(login string, permissionType model.PermissionType) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	types, ok := c.cache[login]
	if ok {
		for i := range types {
			if types[i] == permissionType {
				return ErrPermissionAlreadyExists
			}
		}

		types = append(types, permissionType)
		c.cache[login] = types

		return nil
	}

	c.cache[login] = []model.PermissionType{permissionType}

	return nil
}

func (c *PermissionsCache) GetPermissions(login string) ([]model.PermissionType, error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	perms, ok := c.cache[login]
	if !ok {
		return nil, ErrPermissionNotFound
	}

	return perms, nil
}
