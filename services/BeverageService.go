package services

import (
	"github.com/killingspark/beverages/models"
)

//BeverageService is used for persistence management for Beverages
type BeverageService struct {
	beverages []models.Beverage
}

func MakeBeverageService() BeverageService {
	var bs BeverageService
	var bev = models.Beverage{ID: 0, Name: "Bier", Value: 100}
	bs.beverages = append(bs.beverages, bev)
	return bs
}

func (this *BeverageService) GetBeverages() []models.Beverage {
	return this.beverages
}

func (this *BeverageService) GetBeverage(aID int64) (models.Beverage, bool) {
	if aID >= int64(len(this.beverages)) {
		return models.Beverage{}, false
	}
	return this.beverages[aID], true
}

func (this *BeverageService) NewBeverage(aName string, aValue int) (models.Beverage, bool) {
	var bev models.Beverage
	bev.ID = int64(len(this.beverages))
	bev.Name = aName
	bev.Value = aValue
	this.beverages = append(this.beverages, bev)

	return bev, true
}

func (this *BeverageService) UpdateBeverage(aID int64, aName string, aValue int) (models.Beverage, bool) {
	this.beverages[aID].Name = aName
	this.beverages[aID].Value = aValue
	return this.beverages[aID], true
}

func (this *BeverageService) DeleteBeverage(aID int64) bool {
	this.beverages[aID].Value = -1337
	this.beverages[aID].Name = "UNVALID"
	this.beverages[aID].ID = -1
	return true
}
