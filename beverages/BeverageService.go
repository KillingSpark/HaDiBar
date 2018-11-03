package beverages

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"strconv"

	"github.com/killingspark/HaDiBar/settings"
)

//BeverageService handles the persistence of beverages for us
type BeverageService struct {
	path      string
	beverages map[string]*Beverage
}

//NewBeverageService creates a new Service
func NewBeverageService() *BeverageService {
	bs := &BeverageService{}
	bs.path = settings.S.BeveragePath
	return bs
}

func (service *BeverageService) Load() error {
	if _, err := os.Stat(service.path); os.IsNotExist(err) {
		service.beverages = make(map[string]*Beverage)
		return nil
	}
	jsonFile, err := os.Open(service.path)
	// if we os.Open returns an error then handle it
	if err != nil {
		service.beverages = make(map[string]*Beverage)
		return err
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		service.beverages = make(map[string]*Beverage)
		return err
	}

	if len(byteValue) == 0 { //empty file
		service.beverages = make(map[string]*Beverage)
		return nil
	}

	err = json.Unmarshal([]byte(byteValue), &service.beverages)
	if err != nil {
		service.beverages = make(map[string]*Beverage)
		return err
	}
	return nil
}

func (service *BeverageService) Save() error {
	os.Remove(service.path)
	os.Create(service.path)
	jsonFile, err := os.OpenFile(service.path, os.O_RDWR, 0)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	enc, err := json.Marshal(service.beverages)
	if err != nil {
		return err
	}

	_, err = jsonFile.Write(enc)
	if err != nil {
		return err
	}

	return nil
}

//GetBeverages returns all existing beverages
func (service *BeverageService) GetBeverages(groupID string) ([]*Beverage, error) {
	err := service.Load()
	if err != nil {
		return nil, err
	}
	var bevs []*Beverage
	for _, bev := range service.beverages {
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
	err := service.Load()
	if err != nil {
		return nil, err
	}
	bev, ok := service.beverages[aID]
	if !ok {
		return nil, ErrInvalidID
	}
	if bev.GroupID != groupID {
		return nil, ErrInvalidGroupID
	}

	return bev, nil
}

//NewBeverage creates a new beverage and stores it in the database
func (service *BeverageService) NewBeverage(groupId, aName string, aValue int) (*Beverage, error) {
	err := service.Load()
	if err != nil {
		return nil, err
	}
	bev := &Beverage{ID: strconv.FormatInt(time.Now().UnixNano(), 10), GroupID: groupId, Name: aName, Value: aValue}

	service.beverages[bev.ID] = bev

	err = service.Save()
	if err != nil {
		return nil, err
	}

	return bev, nil
}

//UpdateBeverage updates the data for the identified beverage (eg name and value)
func (service *BeverageService) UpdateBeverage(aID string, aName string, aValue int) (*Beverage, error) {
	err := service.Load()
	if err != nil {
		return nil, err
	}
	bev, ok := service.beverages[aID]
	if !ok {
		return nil, ErrInvalidID
	}
	bev.Name = aName
	bev.Value = aValue

	err = service.Save()
	if err != nil {
		return nil, err
	}

	return bev, nil
}

//DeleteBeverage deletes the identified beverage
func (service *BeverageService) DeleteBeverage(aID string) error {
	err := service.Load()
	if err != nil {
		return err
	}
	_, ok := service.beverages[aID]
	if !ok {
		return ErrInvalidID
	}

	delete(service.beverages, aID)

	err = service.Save()
	if err != nil {
		return err
	}

	return nil
}
