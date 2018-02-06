package accounts

//Entity (s) represent owners of an Account
type Entity struct {
	Name  string `json:"Name"`
	Floor string `json:"Floor"`
	ID    int64  `json:"ID"`
}
