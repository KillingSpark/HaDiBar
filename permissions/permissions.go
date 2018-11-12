package permissions

import (
	"errors"
	"os"

	"github.com/nanobox-io/golang-scribble"
)

type ObjectID string
type UserID string
type PermissionType int8

const (
	Create PermissionType = 0
	Read   PermissionType = 1
	Update PermissionType = 2
	Delete PermissionType = 3
	CRUD   PermissionType = 4

	collection = "permissions"
	resource   = "permissions"
)

var (
	ErrUserIDNotKnow    = errors.New("Userid has no permission set for this object")
	ErrObjectIDNotKnow  = errors.New("Objectid not know to permission")
	ErrPermissionNotSet = errors.New("PermissionType not set for useid,ojectid tuple")
)

type Permissions struct {
	//maps from objectID,userID to the permissiontypes given
	permmap           map[ObjectID](map[UserID](map[PermissionType]bool))
	permRepo          *scribble.Driver
	defaultPermission bool
}

func NewPermissions(path string) (*Permissions, error) {
	perm := &Permissions{}
	var err error
	perm.permRepo, err = scribble.New(path, nil)
	if err != nil {
		return nil, err
	}
	err = perm.permRepo.Read("permissions", "permissions", &perm.permmap)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if perm.permmap == nil {
		perm.permmap = make(map[ObjectID](map[UserID](map[PermissionType]bool)))
	}
	return perm, nil
}

func (p *Permissions) SetPermission(objID, usrID string, permission PermissionType, value bool) error {
	o, ok := p.permmap[ObjectID(objID)]
	if !ok {
		p.permmap[ObjectID(objID)] = make(map[UserID](map[PermissionType]bool))
		o = p.permmap[ObjectID(objID)]
	}
	u, ok := o[UserID(usrID)]
	if !ok {
		o[UserID(usrID)] = make(map[PermissionType]bool)
		u = o[UserID(usrID)]
	}
	u[permission] = value
	err := p.permRepo.Write(collection, resource, p.permmap)
	return err
}

func (p *Permissions) DeletePermission(objID, usrID string, permission PermissionType) error {
	o, ok := p.permmap[ObjectID(objID)]
	if !ok {
		return ErrObjectIDNotKnow
	}
	u, ok := o[UserID(usrID)]
	if !ok {
		return ErrUserIDNotKnow
	}
	delete(u, permission)
	err := p.permRepo.Write(collection, resource, p.permmap)
	return err
}

func (p *Permissions) CheckPermissionAny(objID, usrID string, permissions ...PermissionType) (bool, error) {
	for _, permission := range permissions {
		o, ok := p.permmap[ObjectID(objID)]
		if ok {
			u, ok := o[UserID(usrID)]
			if ok {
				prm, ok := u[permission]
				if ok {
					if prm {
						return true, nil
					}
				} else {
					if p.defaultPermission {
						return true, nil
					} else {
						continue
					}
				}
			} else {
				return p.defaultPermission, ErrUserIDNotKnow
			}
		} else {
			return p.defaultPermission, ErrObjectIDNotKnow
		}
	}
	return false, nil
}
