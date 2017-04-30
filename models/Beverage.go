package models

//Beverage : Model for the system
type Beverage struct {
	ID    int64  `json:"ID"`
	Name  string `json:"Name"`
	Value int    `json:"Value"`
}
