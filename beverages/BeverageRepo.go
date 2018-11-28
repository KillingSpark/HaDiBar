package beverages

import (
	"encoding/json"
	scribble "github.com/nanobox-io/golang-scribble"
)

type BeverageRepo struct {
	db *scribble.Driver
}

var collectionName = "beverages"

func NewBeverageRepo(path string) (*BeverageRepo, error) {
	ar := &BeverageRepo{}
	var err error
	ar.db, err = scribble.New(path, nil)
	if err != nil {
		return nil, err
	}
	return ar, nil
}

func (ar *BeverageRepo) GetAllBeverages() ([]*Beverage, error) {
	list, err := ar.db.ReadAll(collectionName)
	if err != nil {
		return nil, err
	}

	var res []*Beverage
	for _, item := range list {
		bev := &Beverage{}
		err := json.Unmarshal([]byte(item), bev)
		if err != nil {
			continue //skip invalied entries. maybe implement cleanup...
		}
		res = append(res, bev)
	}
	return res, nil
}

func (ar *BeverageRepo) SaveInstance(bev *Beverage) error {
	if err := ar.db.Write(collectionName, bev.ID, bev); err != nil {
		return err
	}
	return nil
}

func (ar *BeverageRepo) GetInstance(bevID string) (*Beverage, error) {
	var bev Beverage
	if err := ar.db.Read(collectionName, bevID, &bev); err != nil {
		return nil, err
	}
	return &bev, nil
}

func (ar *BeverageRepo) DeleteInstance(bevID string) error {
	return ar.db.Delete(collectionName, bevID)
}
