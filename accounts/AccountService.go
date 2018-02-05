package accounts

import (
	"errors"
	"strconv"
)

//AccountService is a service for accessing accounts
type AccountService struct {
	accounts []Account
}

//NewAccountService creates a AccountService and initialzes the Data
func NewAccountService() *AccountService {
	var acs AccountService

	for i := 0; i < 10; i++ {
		acc := Account{ID: int64(i), Owner: Entity{ID: int64(1337 + i), Name: "Moritz" + strconv.Itoa(i)}, Value: 4242}
		acs.accounts = append(acs.accounts, acc)
	}
	return &acs
}

//GetAccounts returns all accounts
func (service *AccountService) GetAccounts() []Account {
	return service.accounts
}

//GetAccount returns the account indentified by accounts/:id
func (service *AccountService) GetAccount(aID int64) Account {
	return service.accounts[aID]
}

//UpdateAccount updates the account with the difference and returns the new account
func (service *AccountService) UpdateAccount(userToken string, aID int64, aValue int) (Account, error) {
	if userToken == "" {
		return Account{}, errors.New("no token provided")
	}
	service.accounts[aID].Value += aValue
	return service.accounts[aID], nil
}
