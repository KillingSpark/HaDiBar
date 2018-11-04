package beverages

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/nanobox-io/golang-scribble"

	"strconv"
)

//BeverageService handles the persistence of beverages for us
type BeverageService struct {
	path    string
	bevRepo *scribble.Driver
}

var collectionName = "beverages"

//NewBeverageService creates a new Service
func NewBeverageService(path string) (*BeverageService, error) {
	bs := &BeverageService{}
	var err error
	bs.bevRepo, err = scribble.New(path, nil)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

//GetBeverages returns all existing beverages
func (service *BeverageService) GetBeverages(groupID string) ([]*Beverage, error) {
	list, err := service.bevRepo.ReadAll(collectionName)
	if err != nil {
		return nil, err
	}
	var bevs []*Beverage
	for _, item := range list {
		var bev *Beverage
		err := json.Unmarshal([]byte(item), bev)
		if err != nil {
			continue
		}
		if bev.GroupID == groupID {
			bevs = append(bevs, bev)
		}
	}
	return bevs, nil
}

var ErrInvalidID = errors.New("ID for beverage is invalid")
var ErrInvalidGroupID = errors.New("ID for beverage is not in your group")

//GetBeverage returns the identified beverage
func (service *BeverageService) GetBeverage(aID, groupID string) (*Beverage, error) {
	var bev *Beverage
	err := service.bevRepo.Read(collectionName, aID, bev)
	if err != nil {
		return nil, ErrInvalidID
	}
	if bev.GroupID != groupID {
		return nil, ErrInvalidGroupID
	}

	return bev, nil
}

//NewBeverage creates a new beverage and stores it in the database
func (service *BeverageService) NewBeverage(groupId, aName string, aValue int) (*Beverage, error) {
	bev := &Beverage{ID: strconv.FormatInt(time.Now().UnixNano(), 10), GroupID: groupId, Name: aName, Value: aValue}

	if err := service.bevRepo.Write(collectionName, bev.ID, bev); err != nil {
		return nil, err
	}

	return bev, nil
}

//UpdateBeverage updates the data for the identified beverage (eg name and value)
func (service *BeverageService) UpdateBeverage(aID string, aName string, aValue int) (*Beverage, error) {
	var bev *Beverage
	err := service.bevRepo.Read(collectionName, aID, bev)
	if err != nil {
		return nil, ErrInvalidID
	}
	bev.Name = aName
	bev.Value = aValue

	err = service.bevRepo.Write(collectionName, aID, bev)
	if err != nil {
		return nil, err
	}

	return bev, nil
}

//DeleteBeverage deletes the identified beverage
func (service *BeverageService) DeleteBeverage(aID string) error {
	var bev *Beverage
	err := service.bevRepo.Read(collectionName, aID, bev)
	if err != nil {
		return ErrInvalidID
	}

	err = service.bevRepo.Delete(collectionName, aID)
	if err != nil {
		return err
	}

	return nil
}
