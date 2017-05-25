package beverages

//Beverage : Model for the system
type Beverage struct {
	ID    string `json:"ID"`
	Name  string `json:"Name"`
	Value int    `json:"Value"`
}
