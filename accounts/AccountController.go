package accounts

import (
	"encoding/json"
	"fmt"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/killingspark/HaDiBar/restapi"
)

//AccountController is the controller for accounts
type AccountController struct {
	service *AccountService
}

//NewAccountController creates a new AccountController and initializes the service
func NewAccountController() *AccountController {
	var acC AccountController
	acC.service = NewAccountService()
	return &acC
}

//GetAccounts gets all existing accounts
func (controller *AccountController) GetAccounts(ctx *gin.Context) {
	enc, _ := json.Marshal(restapi.Response{Status: "OK", Response: controller.service.GetAccounts()})
	fmt.Fprint(ctx.Writer, string(enc))
}

//GetAccount returns the account identified by account/:id
func (controller *AccountController) GetAccount(ctx *gin.Context) {
	strID, _ := ctx.GetQuery("id")
	ID, _ := strconv.Atoi(strID)
	enc, _ := json.Marshal(restapi.Response{Status: "OK", Response: controller.service.GetAccount(int64(ID))})
	fmt.Fprint(ctx.Writer, string(enc))
}

//UpdateAccount updates the value of the account identified by accounts/:id with the form-value "value" as diffenrence
func (controller *AccountController) UpdateAccount(ctx *gin.Context) {
	strID, _ := ctx.GetQuery("id")
	ID, err := strconv.Atoi(strID)

	if err != nil {
		fmt.Fprint(ctx.Writer, "{\"status\":\"ERROR\", \"reponse\":\"id is NaN\"}")
		return
	}
	value, err := strconv.Atoi(ctx.PostForm("value"))
	if err != nil {
		fmt.Fprint(ctx.Writer, "{\"status\":\"ERROR\", \"reponse\":\"value is NaN\"}")
		return
	}
	acc, _ := controller.service.UpdateAccount(int64(ID), value)
	enc, _ := json.Marshal(restapi.Response{Status: "OK", Response: acc})
	fmt.Fprint(ctx.Writer, string(enc))
}
