package beverages

import (
	"database/sql"

	"time"

	"strconv"

	"github.com/killingspark/HaDiBar/logger"
	_ "github.com/mattn/go-sqlite3" //sqlite driver
)

//SQLiteBeverageService handles the persistence of beverages for us
type SQLiteBeverageService struct {
	*sql.DB
}

//NewSQLiteBeverageService creates a new SQLiteService and initialises the database
func NewSQLiteBeverageService() *SQLiteBeverageService {
	db, err := sql.Open("sqlite3", "beverages.db")

	if err != nil {
		logger.Logger.Error("Beverage Database not initialized.")
		db.Close()
	}

	db.Exec("CREATE TABLE beverages (ID int not null, Name char, Value int)")

	var sqlbs SQLiteBeverageService
	sqlbs.DB = db

	return &sqlbs
}

//GetBeverages returns all existing beverages
func (service *SQLiteBeverageService) GetBeverages() []Beverage {
	rows, err := service.Query("SELECT * FROM beverages")
	if err != nil {
		logger.Logger.Error(err.Error())
		return nil
	}

	var bevs []Beverage
	for rows.Next() {
		var bev Beverage
		err := rows.Scan(&bev.ID, &bev.Name, &bev.Value)
		if err != nil {
		}

		bevs = append(bevs, bev)
	}

	return bevs
}

//GetBeverage returns the identified beverage
func (service *SQLiteBeverageService) GetBeverage(aID string) (Beverage, bool) {
	var bev Beverage
	err := service.QueryRow("SELECT * FROM beverages WHERE ID LIKE ?", aID).Scan(&bev.ID, &bev.Name, &bev.Value)
	if err != nil {
		return bev, false
	}

	return bev, true
}

//NewBeverage creates a new beverage and stores it in the database
func (service *SQLiteBeverageService) NewBeverage(aName string, aValue int) (Beverage, bool) {
	bev := Beverage{ID: strconv.FormatInt(time.Now().UnixNano(), 10), Name: aName, Value: aValue}
	_, err := service.Exec("INSERT INTO beverages VALUES (?,?,?)", bev.ID, bev.Name, bev.Value)
	if err != nil {
		return bev, false
	}

	return bev, true
}

//UpdateBeverage updates the data for the identified beverage (eg name and value)
func (service *SQLiteBeverageService) UpdateBeverage(aID string, aName string, aValue int) (Beverage, bool) {
	_, err := service.Exec("UPDATE beverages SET Name = ?, Value = ? WHERE ID = ?", aName, aValue, aID)
	bev := Beverage{ID: aID, Name: aName, Value: aValue}
	if err != nil {
		return bev, false
	}
	return bev, true
}

//DeleteBeverage deletes the identified beverage
func (service *SQLiteBeverageService) DeleteBeverage(aID string) bool {
	_, err := service.Exec("DELETE FROM beverages WHERE ID = ?", aID)
	if err != nil {
		return false
	}
	return true
}
