package accounts

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/killingspark/hadibar/permissions"
	scribble "github.com/nanobox-io/golang-scribble"
)

//AccountService is a service for accessing accounts
type AccountService struct {
	accRepo *scribble.Driver
	perms   *permissions.Permissions
}

var collectionName = "accounts"
var collectionNameTrans = "transactions"

//NewAccountService creates a AccountService and initialzes the Data
func NewAccountService(path string, perms *permissions.Permissions) (*AccountService, error) {
	acs := &AccountService{}
	var err error
	acs.accRepo, err = scribble.New(path, nil)
	if err != nil {
		return nil, err
	}
	acs.perms = perms
	return acs, nil
}

var ErrIDAlreadyTaken = errors.New("AccountID already taken")

func (service *AccountService) Add(new *Account, userID string, perm permissions.PermissionType, perms ...permissions.PermissionType) error {
	service.perms.SetPermission(new.ID, userID, perm, true)
	for _, perm := range perms {
		service.perms.SetPermission(new.ID, userID, perm, true)
	}

	if err := service.accRepo.Write(collectionName, new.ID, new); err != nil {
		return err
	}

	return nil
}

func (service *AccountService) CreateAdd(name, userID string, perm permissions.PermissionType, perms ...permissions.PermissionType) (*Account, error) {
	acc := &Account{}
	acc.ID = strconv.FormatInt(time.Now().UnixNano(), 10)
	acc.Owner.Name = name
	acc.Value = 0

	err := service.Add(acc, userID, perm, perms...)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (service *AccountService) isMainAccount(acc *Account) bool {
	return acc.Owner.Name == "bank"
}

func (service *AccountService) containsMainAccount(accs []*Account) bool {
	for _, acc := range accs {
		if service.isMainAccount(acc) {
			return true
		}
	}
	return false
}

func (service *AccountService) addMainAccount(userID string) (*Account, error) {
	return service.CreateAdd("bank", userID, permissions.Read, permissions.Update)
}

//GetAccounts returns all accounts that are part of this group
func (service *AccountService) GetAccounts(userID string) ([]*Account, error) {
	list, err := service.accRepo.ReadAll(collectionName)
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
		ok, _ := service.perms.CheckPermissionAny(acc.ID, userID, permissions.CRUD, permissions.Read)
		if ok {
			res = append(res, acc)
		}
	}

	if !service.containsMainAccount(res) {
		acc, err := service.addMainAccount(userID)
		if err != nil {
			return nil, err
		}
		res = append(res, acc)
	}

	return res, nil
}

//GetAccount returns the account indentified by accounts/:id
func (service *AccountService) GetAccount(accID, userID string) (*Account, error) {
	ok, err := service.perms.CheckPermissionAny(accID, userID, permissions.CRUD, permissions.Read)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNotOwnerOfObject
	}

	acc := &Account{}
	err = service.accRepo.Read(collectionName, accID, acc)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

//UpdateAccount updates the account with the difference and returns the new account
func (service *AccountService) UpdateAccount(accID, userID string, aDiff int) (*Account, error) {
	ok, err := service.perms.CheckPermissionAny(accID, userID, permissions.CRUD, permissions.Update)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNotOwnerOfObject
	}
	acc, err := service.GetAccount(accID, userID)
	if err != nil {
		return nil, err
	}
	acc.Value += aDiff
	err = service.accRepo.Write(collectionName, accID, acc)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

//UpdateAccount updates the account with the difference and returns the new account
func (service *AccountService) Transaction(SourceID, TargetID, userID string, amount int) error {
	if SourceID != "0" { //0 is reserved for infusions from outside the system
		ok, err := service.perms.CheckPermissionAny(SourceID, userID, permissions.CRUD, permissions.Update)
		if err != nil {
			return err
		}
		if !ok {
			return ErrNotOwnerOfObject
		}
		source, err := service.GetAccount(SourceID, userID)
		if err != nil {
			return err
		}
		source.Value -= amount
		err = service.accRepo.Write(collectionName, source.ID, source)
		if err != nil {
			return err
		}
	}
	ok, err := service.perms.CheckPermissionAny(TargetID, userID, permissions.CRUD, permissions.Update)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotOwnerOfObject
	}
	target, err := service.GetAccount(TargetID, userID)
	if err != nil {
		return err
	}

	target.Value += amount
	err = service.accRepo.Write(collectionName, target.ID, target)
	if err != nil {
		return err
	}
	trans := &Transaction{}
	trans.SourceID = SourceID
	trans.TargetID = TargetID
	trans.Timestamp = time.Now()
	trans.Amount = amount
	trans.ID = strconv.Itoa(trans.Timestamp.Nanosecond())
	err = service.accRepo.Write(collectionNameTrans, trans.ID, trans)
	service.perms.SetPermission(trans.ID, userID, permissions.CRUD, true)
	if err != nil {
		return err
	}
	return nil
}

func (service *AccountService) GetTransactions(accID, userID string) ([]*Transaction, error) {
	list, err := service.accRepo.ReadAll(collectionNameTrans)
	if err != nil {
		return nil, err
	}
	res := make([]*Transaction, 0)
	for _, item := range list {
		tx := &Transaction{}
		err := json.Unmarshal([]byte(item), tx)
		if err != nil {
			continue //skip invalid entries
		}
		if ok, err := service.perms.CheckPermissionAny(tx.ID, userID, permissions.Read, permissions.CRUD); ok {
			if err != nil {
				return nil, err
			}
			if accID == "" || (accID == tx.SourceID || accID == tx.TargetID) {
				res = append(res, tx)
			}
		}
	}
	return res, nil
}

var ErrNotOwnerOfObject = errors.New("This User is not an owner of this account")

func (service *AccountService) GivePermissionToUser(accID, ownerID, newOwnerID string, perm permissions.PermissionType) error {
	ok, err := service.perms.CheckPermissionAny(accID, ownerID, permissions.CRUD, permissions.Read)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotOwnerOfObject
	}

	return service.perms.SetPermission(accID, newOwnerID, perm, true)
}

//UpdateAccount updates the account with the difference and returns the new account
func (service *AccountService) DeleteAccount(accID, userID string) error {
	ok, err := service.perms.CheckPermissionAny(accID, userID, permissions.Delete, permissions.CRUD)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotOwnerOfObject
	}
	err = service.accRepo.Delete(collectionName, accID)
	if err != nil {
		return err
	}
	return nil
}
