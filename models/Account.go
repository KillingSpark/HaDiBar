package models

//Account : represents the Accouts of Floors/Roomates/etc
type Account struct {
	Owner Entity `json:"Owner"`
	Value int    `json:"Value"`
	ID    int64  `json:"ID"`
}
