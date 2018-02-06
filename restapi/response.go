package restapi

import "encoding/json"

//Response represents the strcuture for  response
type Response struct {
	Status   string      `json:"status"`
	Response interface{} `json:"response"`
}

//NewErrorResponse returns a new response with "ERROR" status and the object as response
func NewErrorResponse(resp interface{}) *Response {
	return &Response{Status: "ERROR", Response: resp}
}

//NewOkResponse returns a new response with "OK" status and the object as response
func NewOkResponse(resp interface{}) *Response {
	return &Response{Status: "OK", Response: resp}
}

//NewNosesResponse returns a new response with "OK" status and the object as response
func NewNosesResponse(resp interface{}) *Response {
	return &Response{Status: "NOSES", Response: resp}
}

//NewNoauthResponse returns a new response with "OK" status and the object as response
func NewNoauthResponse(resp interface{}) *Response {
	return &Response{Status: "NOAUTH", Response: resp}
}

func (resp *Response) Marshal() ([]byte, error) {
	enc, err := json.Marshal(resp)
	return enc, err
}
