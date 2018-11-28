package accounts

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"path"
)

var globdb *bolt.DB

type AccountRepo struct {
	db *bolt.DB
}

var collectionName = "Accounts"
var collectionNameTrans = "Transactions"

func NewAccountRepo(dir string) (*AccountRepo, error) {
	ar := &AccountRepo{}
	var err error

	if globdb == nil {
		globdb, err = bolt.Open(path.Join(dir, collectionName+".bolt"), 0600, nil)
		if err != nil {
			return nil, err
		}
		globdb.Update(func(tx *bolt.Tx) error {
			tx.CreateBucket([]byte(collectionName))
			tx.CreateBucket([]byte(collectionNameTrans))
			return nil
		})
	}
	ar.db = globdb
	return ar, nil
}

func (ar *AccountRepo) GetAllAccounts() ([]*Account, error) {
	var res []*Account
	err := ar.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collectionName))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			acc := &Account{}
			err := json.Unmarshal([]byte(v), acc)
			if err != nil {
				continue //skip invalied entries. maybe implement cleanup...
			}
			res = append(res, acc)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ar *AccountRepo) SaveInstance(acc *Account) error {
	err := ar.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collectionName))
		marshed, err := json.Marshal(acc)
		if err != nil {
			return err
		}
		return b.Put([]byte(acc.ID), marshed)
	})
	return err
}

var ErrAccountDoesNotExist = errors.New("Account with this id does not exist")

func (ar *AccountRepo) GetInstance(accID string) (*Account, error) {
	var acc Account
	err := ar.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collectionName))
		marshed := b.Get([]byte(accID))
		if marshed == nil {
			return ErrAccountDoesNotExist
		}
		return json.Unmarshal(marshed, &acc)
	})
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func (ar *AccountRepo) DeleteInstance(accID string) error {
	err := ar.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collectionName))
		b.Delete([]byte(accID))
		return nil
	})
	return err
}

func (ar *AccountRepo) SaveTransaction(trans *Transaction) error {
	err := ar.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collectionNameTrans))
		marshed, err := json.Marshal(trans)
		if err != nil {
			return err
		}
		return b.Put([]byte(trans.ID), marshed)
	})
	return err
}

func (ar *AccountRepo) GetTransactions() ([]*Transaction, error) {
	var res []*Transaction
	err := ar.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collectionNameTrans))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			trans := &Transaction{}
			err := json.Unmarshal([]byte(v), trans)
			if err != nil {
				continue //skip invalied entries. maybe implement cleanup...
			}
			res = append(res, trans)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
