package services

import "github.com/killingspark/beverages/models"

type AccountService struct {
	accounts []models.Account
}

func MakeAccountService() AccountService {
	var acs AccountService
	var acc models.Account
	acc.ID = 0
	acc.Owner = models.Entity{ID: 1337, Name: "Moritz"}
	acc.Value = 4242
	acs.accounts = append(acs.accounts, acc)
	return acs
}

func (this *AccountService) GetAccounts() []models.Account {
	return this.accounts
}

func (this *AccountService) GetAccount(aID int64) models.Account {
	return this.accounts[aID]
}

func (this *AccountService) UpdateAccount(aID int64, aValue int) (models.Account, bool) {
	this.accounts[aID].Value += aValue
	return this.accounts[aID], true
}
