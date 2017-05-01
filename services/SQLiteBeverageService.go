package services

import (
	"database/sql"

	//sqlite driver

	"time"

	"github.com/killingspark/beverages/models"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteBeverageService struct {
	*sql.DB
}

func MakeSQLiteBeverageService() SQLiteBeverageService {
	db, err := sql.Open("sqlite3", "beverages.db")

	if err != nil {
		print("Beverage Database not initialized.")
		db.Close()
	}

	db.Exec("CREATE TABLE beverages (ID int not null, Name char, Value int)")

	var sqlbs SQLiteBeverageService
	sqlbs.DB = db

	return sqlbs
}

func (service *SQLiteBeverageService) GetBeverages() []models.Beverage {
	rows, err := service.Query("SELECT * FROM beverages")
	if err != nil {
		print(err.Error())
		return nil
	}

	var bevs []models.Beverage
	for rows.Next() {
		var bev models.Beverage
		err := rows.Scan(&bev.ID, &bev.Name, &bev.Value)
		if err != nil {
		}

		bevs = append(bevs, bev)
	}

	return bevs
}

func (service *SQLiteBeverageService) GetBeverage(aID int64) (models.Beverage, bool) {
	var bev models.Beverage
	err := service.QueryRow("").Scan(&bev.ID, &bev.Name, &bev.Value)
	if err != nil {
		return bev, false
	}

	return bev, true
}

func (service *SQLiteBeverageService) NewBeverage(aName string, aValue int) (models.Beverage, bool) {
	bev := models.Beverage{ID: time.Now().UnixNano(), Name: aName, Value: aValue}
	_, err := service.Exec("INSERT INTO beverages VALUES (?,?,?)", bev.ID, bev.Name, bev.Value)
	if err != nil {
		return bev, false
	}

	return bev, true
}

func (service *SQLiteBeverageService) UpdateBeverage(aID int64, aName string, aValue int) (models.Beverage, bool) {
	_, err := service.Exec("UPDATE beverages SET Name = ?, Value = ? WHERE ID = ?", aName, aValue, aID)
	bev := models.Beverage{ID: aID, Name: aName, Value: aValue}
	if err != nil {
		return bev, false
	}
	return bev, true
}

func (service *SQLiteBeverageService) DeleteBeverage(aID int64) bool {
	_, err := service.Exec("DELETE FROM beverages WHERE ID = ?", aID)
	if err != nil {
		return false
	}
	return true
}
