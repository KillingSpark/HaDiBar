package beverages

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"path"
)

var globdb *bolt.DB

type BeverageRepo struct {
	db *bolt.DB
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

func (ar *BeverageRepo) GetAllBeverages() ([]*Beverage, error) {
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
	err := ar.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		b.Delete([]byte(bevID))
		return nil
	})
	return err
}
