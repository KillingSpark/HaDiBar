package accounts

import (
	"encoding/json"
	scribble "github.com/nanobox-io/golang-scribble"
)

type AccountRepo struct {
	db *scribble.Driver
}

var collectionName = "accounts"
var collectionNameTrans = "transactions"

func NewAccountRepo(path string) (*AccountRepo, error) {
	ar := &AccountRepo{}
	var err error
	ar.db, err = scribble.New(path, nil)
	if err != nil {
		return nil, err
	}
	return ar, nil
}

func (ar *AccountRepo) GetAllAccounts() ([]*Account, error) {
	list, err := ar.db.ReadAll(collectionName)
	if err != nil {
		return nil, err
	}

	var res []*Account
	for _, item := range list {
		acc := &Account{}
		err := json.Unmarshal([]byte(item), acc)
		if err != nil {
			continue //skip invalied entries. maybe implement cleanup...
		}
		res = append(res, acc)
	}
	return res, nil
}

func (ar *AccountRepo) SaveInstance(acc *Account) error {
	if err := ar.db.Write(collectionName, acc.ID, acc); err != nil {
		return err
	}
	return nil
}

func (ar *AccountRepo) GetInstance(accID string) (*Account, error) {
	var acc Account
	if err := ar.db.Read(collectionName, accID, &acc); err != nil {
		return nil, err
	}
	return &acc, nil
}

func (ar *AccountRepo) DeleteInstance(accID string) error {
	return ar.db.Delete(collectionName, accID)
}

func (ar *AccountRepo) SaveTransaction(tx *Transaction) error {
	return ar.db.Write(collectionNameTrans, tx.ID, tx)
}

func (ar *AccountRepo) GetTransactions() ([]*Transaction, error) {
	list, err := ar.db.ReadAll(collectionNameTrans)
	if err != nil {
		return nil, err
	}

	var res []*Transaction
	for _, item := range list {
		acc := &Transaction{}
		err := json.Unmarshal([]byte(item), acc)
		if err != nil {
			continue //skip invalied entries. maybe implement cleanup...
		}
		res = append(res, acc)
	}
	return res, nil
}
