package restapi

//Response represents the strcuture for  response
type Response struct {
	Status   string      `json:"status"`
	Response interface{} `json:"response"`
}
