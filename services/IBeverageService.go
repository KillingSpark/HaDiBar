package services

import "github.com/killingspark/HaDiBar/models"

//IBeverageService is the interface needed to interact with beverages
type IBeverageService interface {
	GetBeverages() []models.Beverage
	GetBeverage(aID string) (models.Beverage, bool)
	NewBeverage(aName string, aValue int) (models.Beverage, bool)
	UpdateBeverage(aID string, aName string, aValue int) (models.Beverage, bool)
	DeleteBeverage(aID string) bool
}
