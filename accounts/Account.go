package accounts

import "time"

//Account : represents the Accouts of Floors/Roomates/etc
type Account struct {
	Value int          `json:"Value"`
	ID    string       `json:"ID"`
	Owner AccountOwner `json:"Owner"`
}

type AccountOwner struct {
	Name string `json:"Name"`
}

type Transaction struct {
	SourceID  string
	TargetID  string
	Amount    int
	Timestamp time.Time
	ID        string
}
