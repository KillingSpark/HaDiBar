package authStuff

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"os"
	"path"
)

var globdb *bolt.DB

type UserRepo struct {
	db *bolt.DB
}

var bucketName = "users"

func NewUserRepo(dir string) (*UserRepo, error) {
	ur := &UserRepo{}
	var err error

	if globdb == nil {
		globdb, err = bolt.Open(path.Join(dir, bucketName+".bolt"), 0600, nil)
		if err != nil {
			return nil, err
		}
		globdb.Update(func(tx *bolt.Tx) error {
			tx.CreateBucket([]byte(bucketName))
			return nil
		})
	}
	ur.db = globdb
	return ur, nil
}

func (ur *UserRepo) BackupTo(bkpDest string) error {
	f, err := os.OpenFile(bkpDest, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	err = ur.db.View(func(tx *bolt.Tx) error {
		_, err = tx.WriteTo(f)
		return err
	})
	return err
}

func (ur *UserRepo) GetAllUsers() ([]*LoginInfo, error) {
	var res []*LoginInfo
	err := ur.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			info := &LoginInfo{}
			err := json.Unmarshal([]byte(v), info)
			if err != nil {
				continue //skip invalied entries. maybe implement cleanup...
			}
			res = append(res, info)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ur *UserRepo) SaveInstance(info *LoginInfo) error {
	err := ur.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		marshed, err := json.Marshal(info)
		if err != nil {
			return err
		}
		return b.Put([]byte(info.Name), marshed)
	})
	return err
}

var ErrUserDoesNotExist = errors.New("User with this Name does not exist")

func (ur *UserRepo) GetInstance(infoName string) (*LoginInfo, error) {
	var info LoginInfo
	err := ur.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		marshed := b.Get([]byte(infoName))
		if marshed == nil {
			return ErrUserDoesNotExist
		}
		return json.Unmarshal(marshed, &info)
	})
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (ur *UserRepo) DeleteInstance(infoName string) error {
	err := ur.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		b.Delete([]byte(infoName))
		return nil
	})
	return err
}
