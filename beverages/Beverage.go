package beverages

//Beverage : Model for the system
type Beverage struct {
	ID       string   `json:"ID"`
	GroupIDs []string `json:"GroupIDs"`
	Name     string   `json:"Name"`
	Value    int      `json:"Value"`
}
