package accounts

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/nanobox-io/golang-scribble"
)

//AccountService is a service for accessing accounts
type AccountService struct {
	accounts map[string]*Account
	accRepo  *scribble.Driver
}

var collectionName = "accounts"

//NewAccountService creates a AccountService and initialzes the Data
func NewAccountService(path string) (*AccountService, error) {
	acs := &AccountService{}
	var err error
	acs.accRepo, err = scribble.New(path, nil)
	if err != nil {
		return nil, err
	}
	return acs, nil
}

func dummyAccs() []*Account {
	accs := make([]*Account, 0)
	for i := 0; i < 10; i++ {
		acc := &Account{ID: strconv.Itoa(i), Value: 4242, Groups: []*AccountGroup{&AccountGroup{GroupID: "M6"}}, Owner: AccountOwner{Name: "Moritz" + strconv.Itoa(i)}}
		accs = append(accs, acc)
	}
	for i := 0; i < 10; i++ {
		acc := &Account{ID: strconv.Itoa(i), Value: 4242, Groups: []*AccountGroup{&AccountGroup{GroupID: "M5"}}, Owner: AccountOwner{Name: "Paul" + strconv.Itoa(i)}}
		accs = append(accs, acc)
	}
	return accs
}

var ErrIDAlreadyTaken = errors.New("AccountID already taken")

func (service *AccountService) Add(new *Account) error {
	acc, err := service.GetAccount(new.ID)
	if err != nil {
		return err
	}
	if acc != nil {
		return ErrIDAlreadyTaken
	}

	if err := service.accRepo.Write(collectionName, new.ID, new); err != nil {
		return err
	}

	return nil
}

func (service *AccountService) CreateAdd(name, groupID string) (*Account, error) {
	acc := &Account{}
	acc.ID = strconv.FormatInt(time.Now().UnixNano(), 10)
	acc.Owner.Name = name
	acc.Groups = []*AccountGroup{&AccountGroup{GroupID: groupID}}
	acc.Value = 0

	if err := service.accRepo.Write(collectionName, acc.ID, acc); err != nil {
		return nil, err
	}

	return nil, nil
}

//GetAccounts returns all accounts that are part of this group
func (service *AccountService) GetAccounts(groupID string) ([]*Account, error) {
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
		for _, group := range acc.Groups {
			if group.GroupID == groupID {
				res = append(res, acc)
				break
			}
		}
	}
	return res, nil
}

//GetAccount returns the account indentified by accounts/:id
func (service *AccountService) GetAccount(aID string) (*Account, error) {
	acc := &Account{}
	err := service.accRepo.Read(collectionName, aID, acc)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

//UpdateAccount updates the account with the difference and returns the new account
func (service *AccountService) UpdateAccount(aID string, aDiff int) (*Account, error) {
	acc, err := service.GetAccount(aID)
	if err != nil {
		return nil, err
	}
	acc.Value += aDiff
	err = service.accRepo.Write(collectionName, aID, acc)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

var ErrNotOwnerOfObject = errors.New("This User is not an owner of this account")

func (service *AccountService) AddAccountToGroup(aID, groupID, aNewGroup string) (*Account, error) {
	acc, err := service.GetAccount(aID)
	if err != nil {
		return nil, err
	}
	contains := false
	for _, group := range acc.Groups {
		if group.GroupID == groupID {
			contains = true
			break
		}
	}
	if !contains {
		return nil, ErrNotOwnerOfObject
	}
	acc.Groups = append(acc.Groups, &AccountGroup{GroupID: aNewGroup})
	err = service.accRepo.Write(collectionName, aID, acc)
	if err != nil {
		return nil, err
	}
	return acc, nil
}
