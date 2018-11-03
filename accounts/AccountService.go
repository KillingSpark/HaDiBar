package accounts

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/killingspark/HaDiBar/settings"
)

//AccountService is a service for accessing accounts
type AccountService struct {
	accounts []*Account
	path     string
}

//NewAccountService creates a AccountService and initialzes the Data
func NewAccountService() *AccountService {
	var acs AccountService
	acs.path = os.ExpandEnv(settings.S.AccountPath)
	return &acs
}

func dummyAccs() []*Account {
	accs := make([]*Account, 0)
	for i := 0; i < 10; i++ {
		acc := &Account{ID: int64(i), Value: 4242, Group: AccountGroup{GroupID: "M6"}, Owner: AccountOwner{Name: "Moritz" + strconv.Itoa(i)}}
		accs = append(accs, acc)
	}
	for i := 0; i < 10; i++ {
		acc := &Account{ID: int64(i), Value: 4242, Group: AccountGroup{GroupID: "M5"}, Owner: AccountOwner{Name: "Paul" + strconv.Itoa(i)}}
		accs = append(accs, acc)
	}
	return accs
}

var ErrIDAlreadyTaken = errors.New("AccountID already taken")

func (service *AccountService) Add(new *Account) error {
	for _, acc := range service.accounts {
		if acc.ID == new.ID {
			return ErrIDAlreadyTaken
		}
	}
	service.accounts = append(service.accounts, new)
	return nil
}

func (service *AccountService) Load() error {
	jsonFile, err := os.Open(service.path)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	if len(byteValue) == 0 { //empty file
		service.accounts = make([]*Account, 0)
		return nil
	}

	err = json.Unmarshal([]byte(byteValue), &service.accounts)
	if err != nil {
		service.accounts = make([]*Account, 0)
		return err
	}
	return nil
}

func (service *AccountService) Save() error {
	jsonFile, err := os.OpenFile(service.path, os.O_RDWR, 0)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	enc, err := json.Marshal(service.accounts)
	if err != nil {
		return err
	}

	_, err = jsonFile.Write(enc)
	if err != nil {
		return err
	}

	return nil
}

//GetAccounts returns all accounts that are part of this group
func (service *AccountService) GetAccounts(groupID string) []*Account {
	err := service.Load()
	if err != nil {
		return nil
	}
	var res []*Account
	for _, acc := range service.accounts {
		if acc.Group.GroupID == groupID {
			res = append(res, acc)
		}
	}
	return res
}

//GetAccount returns the account indentified by accounts/:id
func (service *AccountService) GetAccount(aID int64) (*Account, error) {
	err := service.Load()
	if err != nil {
		return nil, err
	}
	return service.accounts[aID], nil
}

//UpdateAccount updates the account with the difference and returns the new account
func (service *AccountService) UpdateAccount(aID int64, aValue int) (*Account, error) {
	service.accounts[aID].Value += aValue
	err := service.Save()
	if err != nil {
		return nil, err
	}
	return service.accounts[aID], nil
}
