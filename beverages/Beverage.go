package beverages

//Beverage : Model for the system
type Beverage struct {
	ID      string `json:"ID"`
	GroupID string `json:"GroupID"`
	Name    string `json:"Name"`
	Value   int    `json:"Value"`
}
