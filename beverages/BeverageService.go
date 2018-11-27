package beverages

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/killingspark/hadibar/permissions"
	scribble "github.com/nanobox-io/golang-scribble"

	"strconv"
)

//BeverageService handles the persistence of beverages for us
type BeverageService struct {
	path    string
	bevRepo *scribble.Driver
	perms   *permissions.Permissions
}

var collectionName = "beverages"

var ErrInvalidID = errors.New("ID for beverage is invalid")
var ErrInvalidGroupID = errors.New("ID for beverage is not in your group")
var ErrNoPermission = errors.New("No permission for this action")

//NewBeverageService creates a new Service
func NewBeverageService(path string, perms *permissions.Permissions) (*BeverageService, error) {
	bs := &BeverageService{}
	var err error
	bs.bevRepo, err = scribble.New(path, nil)
	if err != nil {
		return nil, err
	}
	bs.perms = perms
	return bs, nil
}

//GetBeverages returns all existing beverages
func (service *BeverageService) GetBeverages(userID string) ([]*Beverage, error) {
	list, err := service.bevRepo.ReadAll(collectionName)
	if err != nil {
		return nil, err
	}
	var bevs []*Beverage
	for _, item := range list {
		bev := &Beverage{}
		err := json.Unmarshal([]byte(item), bev)
		if err != nil {
			continue
		}
		if ok, _ := service.perms.CheckPermissionAny(bev.ID, userID, permissions.Read, permissions.CRUD); ok {
			bevs = append(bevs, bev)
		}
	}
	return bevs, nil
}

//GetBeverage returns the identified beverage
func (service *BeverageService) GetBeverage(bevID, userID string) (*Beverage, error) {
	bev := &Beverage{}
	err := service.bevRepo.Read(collectionName, bevID, bev)
	if err != nil {
		return nil, ErrInvalidID
	}
	ok, err := service.perms.CheckPermissionAny(bev.ID, userID, permissions.Read, permissions.CRUD)
	if ok {
		return bev, nil
	}

	return nil, err
}

//NewBeverage creates a new beverage and stores it in the database
func (service *BeverageService) NewBeverage(userID, aName string, aValue, aAvailable int) (*Beverage, error) {
	bev := &Beverage{ID: strconv.FormatInt(time.Now().UnixNano(), 10), Name: aName, Value: aValue, Available: aAvailable}

	service.perms.SetPermission(bev.ID, userID, permissions.CRUD, true)

	if err := service.bevRepo.Write(collectionName, bev.ID, bev); err != nil {
		return nil, err
	}

	return bev, nil
}

//UpdateBeverage updates the data for the identified beverage (eg name and value)
func (service *BeverageService) UpdateBeverage(bevID, userID, aName string, aValue, aAvailable int) (*Beverage, error) {

	ok, err := service.perms.CheckPermissionAny(bevID, userID, permissions.Delete, permissions.CRUD)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNoPermission
	}

	bev, err := service.GetBeverage(bevID, userID)
	if err != nil {
		return nil, err
	}
	bev.Name = aName
	bev.Value = aValue
	bev.Available = aAvailable

	err = service.bevRepo.Write(collectionName, bevID, bev)
	if err != nil {
		return nil, err
	}

	return bev, nil
}

//DeleteBeverage deletes the identified beverage
func (service *BeverageService) DeleteBeverage(bevID, userID string) error {

	ok, err := service.perms.CheckPermissionAny(bevID, userID, permissions.Delete, permissions.CRUD)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNoPermission
	}
	err = service.bevRepo.Delete(collectionName, bevID)
	if err != nil {
		return err
	}

	return nil
}

//GivePermissionToUser gives the newOwner the permissions
func (service *BeverageService) GivePermissionToUser(bevID, ownerID, newOwnerID string, perm permissions.PermissionType) error {
	ok, err := service.perms.CheckPermissionAny(bevID, ownerID, permissions.Update, permissions.CRUD)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNoPermission
	}
	return service.perms.SetPermission(bevID, newOwnerID, perm, true)
}
