package permissions

import (
	"errors"
	"os"
	"path"
	"sync"

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
	db                *bolt.DB
	defaultPermission bool
	Lock              sync.RWMutex
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

func (p *Permissions) BackupTo(bkpDest string) error {
	f, err := os.OpenFile(bkpDest, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	err = p.db.View(func(tx *bolt.Tx) error {
		_, err = tx.WriteTo(f)
		return err
	})
	return err
}

func (p *Permissions) RemoveUsersPermissions(usrID string) error {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	err := p.db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(usrID))
		return err
	})
	return err
}

func (p *Permissions) GetAllAsMap() (map[string](map[string](map[PermissionType]bool)), error) {
	p.Lock.RLock()
	defer p.Lock.RUnlock()

	res := make(map[string](map[string](map[PermissionType]bool)))
	err := p.db.View(func(tx *bolt.Tx) error {
		usrc := tx.Cursor()
		for usr, _ := usrc.First(); usr != nil; usr, _ = usrc.Next() {
			res[string(usr)] = make(map[string](map[PermissionType]bool))
			usrb := tx.Bucket(usr)
			objc := usrb.Cursor()
			for obj, _ := objc.First(); obj != nil; obj, _ = objc.Next() {
				res[string(usr)][string(obj)] = make(map[PermissionType]bool)
				objb := usrb.Bucket(obj)
				permc := objb.Cursor()
				for p, v := permc.First(); p != nil; p, v = permc.Next() {
					res[string(usr)][string(obj)][PermissionType(p[0])] = v[0] == 1
				}
			}
		}
		return nil
	})
	return res, err
}

func (p *Permissions) SetPermission(objID, usrID string, permission PermissionType, value bool) error {
	p.Lock.Lock()
	defer p.Lock.Unlock()

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
	p.Lock.Lock()
	defer p.Lock.Unlock()

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
	p.Lock.RLock()
	defer p.Lock.RUnlock()
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
