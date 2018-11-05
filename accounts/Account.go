package accounts

//Account : represents the Accouts of Floors/Roomates/etc
type Account struct {
	Value  int             `json:"Value"`
	ID     string          `json:"ID"`
	Owner  AccountOwner    `json:"Owner"`
	Groups []*AccountGroup `json:"Groups"`
}

type AccountGroup struct {
	GroupID string `json:"GroupID"`
}

type AccountOwner struct {
	Name string `json:"Name"`
}
