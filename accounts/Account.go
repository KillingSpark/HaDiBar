package accounts

//Account : represents the Accouts of Floors/Roomates/etc
type Account struct {
	Value int          `json:"Value"`
	ID    int64        `json:"ID"`
	Owner AccountOwner `json:"Owner"`
	Group AccountGroup `json:"Group"`
}

type AccountGroup struct {
	GroupID string `json:"GroupID"`
}

type AccountOwner struct {
	Name string `json:"Name"`
}
