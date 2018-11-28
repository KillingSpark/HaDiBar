package permissions

import (
	"errors"
	"path"

	"github.com/boltdb/bolt"
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
	db                *bolt.DB
	defaultPermission bool
}

var globdb *bolt.DB

func NewPermissions(dir string) (*Permissions, error) {
	perm := &Permissions{}
	var err error
	if globdb == nil {
		globdb, err = bolt.Open(path.Join(dir, "permissions.bolt"), 0600, nil)
		if err != nil {
			return nil, err
		}
	}
	perm.db = globdb
	return perm, nil
}

func (p *Permissions) SetPermission(objID, usrID string, permission PermissionType, value bool) error {
	err := p.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(usrID))
		if err != nil {
			return err
		}
		o, err := b.CreateBucketIfNotExists([]byte(objID))
		if err != nil {
			return err
		}
		if value {
			err = o.Put([]byte{byte(permission)}, []byte{byte(1)})
		} else {
			err = o.Put([]byte{byte(permission)}, []byte{byte(0)})
		}
		return err
	})
	return err
}

func (p *Permissions) DeletePermission(objID, usrID string, permission PermissionType) error {
	err := p.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(usrID))
		if err != nil {
			return err
		}
		o, err := b.CreateBucketIfNotExists([]byte(objID))
		if err != nil {
			return err
		}
		err = o.Delete([]byte{byte(permission)})
		return err
	})
	return err
}

func (p *Permissions) CheckPermissionAny(objID, usrID string, permissions ...PermissionType) (bool, error) {
	result := false
	err := p.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(usrID))
		if err != nil {
			return err
		}
		o, err := b.CreateBucketIfNotExists([]byte(objID))
		if err != nil {
			return err
		}

		for _, permission := range permissions {
			val := o.Get([]byte{byte(permission)})
			if len(val) < 1 {
				continue
			}
			if val[0] != 0 {
				result = true
				return nil
			}
			if val[0] == 0 {
				result = false
				return nil
			}
		}
		result = p.defaultPermission
		return err
	})
	return result, err
}
