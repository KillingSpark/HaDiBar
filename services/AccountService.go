package services

import "github.com/killingspark/HaDiBar/models"

//AccountService is a service for accessing accounts
type AccountService struct {
	accounts []models.Account
}

//MakeAccountService creates a AccountService and initialzes the Data
func MakeAccountService() AccountService {
	var acs AccountService
	var acc models.Account
	acc.ID = 0
	acc.Owner = models.Entity{ID: 1337, Name: "Moritz"}
	acc.Value = 4242
	acs.accounts = append(acs.accounts, acc)
	return acs
}

//GetAccounts returns all accounts
func (service *AccountService) GetAccounts() []models.Account {
	return service.accounts
}

//GetAccount returns the account indentified by accounts/:id
func (service *AccountService) GetAccount(aID int64) models.Account {
	return service.accounts[aID]
}

//UpdateAccount updates the account with the difference and returns the new account
func (service *AccountService) UpdateAccount(aID int64, aValue int) (models.Account, bool) {
	service.accounts[aID].Value += aValue
	return service.accounts[aID], true
}
