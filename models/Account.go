package models

//Account : represents the Accouts of Floors/Roomates/etc
type Account struct {
	owner Entity
	value int
	ID    int64
}
