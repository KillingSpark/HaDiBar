package services

import "github.com/killingspark/beverages/models"

type IBeverageService interface {
	GetBeverages() []models.Beverage
	GetBeverage(aID int64) (models.Beverage, bool)
	NewBeverage(aName string, aValue int) (models.Beverage, bool)
	UpdateBeverage(aID int64, aName string, aValue int) (models.Beverage, bool)
	DeleteBeverage(aID int64) bool
}
