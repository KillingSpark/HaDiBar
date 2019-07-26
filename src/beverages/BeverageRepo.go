package beverages

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"os"
	"path"
	"sync"
)

var globdb *bolt.DB

type BeverageRepo struct {
	db   *bolt.DB
	Lock sync.RWMutex
}

var bucketName = "beverages"

func NewBeverageRepo(dir string) (*BeverageRepo, error) {
	ar := &BeverageRepo{}
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
	ar.db = globdb
	return ar, nil
}

func (br *BeverageRepo) BackupTo(bkpDest string) error {
	f, err := os.OpenFile(bkpDest, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	err = br.db.View(func(tx *bolt.Tx) error {
		_, err = tx.WriteTo(f)
		return err
	})
	return err
}

func (ar *BeverageRepo) GetAllBeverages() ([]*Beverage, error) {
	ar.Lock.RLock()
	defer ar.Lock.RUnlock()

	var res []*Beverage
	err := ar.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			bev := &Beverage{}
			err := json.Unmarshal([]byte(v), bev)
			if err != nil {
				continue //skip invalied entries. maybe implement cleanup...
			}
			res = append(res, bev)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ar *BeverageRepo) SaveInstance(bev *Beverage) error {
	ar.Lock.Lock()
	defer ar.Lock.Unlock()

	err := ar.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		marshed, err := json.Marshal(bev)
		if err != nil {
			return err
		}
		return b.Put([]byte(bev.ID), marshed)
	})
	return err
}

var ErrBeverageDoesNotExist = errors.New("Beverage with this id does not exist")

func (ar *BeverageRepo) GetInstance(bevID string) (*Beverage, error) {
	ar.Lock.RLock()
	defer ar.Lock.RUnlock()

	var bev Beverage
	err := ar.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		marshed := b.Get([]byte(bevID))
		if marshed == nil {
			return ErrBeverageDoesNotExist
		}
		return json.Unmarshal(marshed, &bev)
	})
	if err != nil {
		return nil, err
	}
	return &bev, nil
}

func (ar *BeverageRepo) DeleteInstance(bevID string) error {
	ar.Lock.Lock()
	defer ar.Lock.Unlock()

	err := ar.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		b.Delete([]byte(bevID))
		return nil
	})
	return err
}
