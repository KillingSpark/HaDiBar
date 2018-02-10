package accounts

import (
	"errors"
	"strconv"

	"github.com/killingspark/HaDiBar/authStuff"
)

//AccountService is a service for accessing accounts
type AccountService struct {
	accounts []Account
}

//NewAccountService creates a AccountService and initialzes the Data
func NewAccountService() *AccountService {
	var acs AccountService

	for i := 0; i < 10; i++ {
		acc := Account{ID: int64(i), Value: 4242, Owner: AccountOwner{Name: "Moritz" + strconv.Itoa(i)}}
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
func (service *AccountService) UpdateAccount(logininfo authStuff.LoginInfo, aID int64, aValue int) (Account, error) {
	if !logininfo.LoggedIn {
		return Account{}, errors.New("not logged in")
	}
	service.accounts[aID].Value += aValue
	return service.accounts[aID], nil
}
