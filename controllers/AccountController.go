package controllers

import (
	"encoding/json"
	"fmt"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/killingspark/HaDiBar/services"
)

//AccountController is the controller for accounts
type AccountController struct {
	service services.AccountService
}

//MakeAccountController creates a new AccountController and initializes the service
func MakeAccountController() AccountController {
	var acC AccountController
	acC.service = services.MakeAccountService()
	return acC
}

//GetAccounts gets all existing accounts
func (controller *AccountController) GetAccounts(ctx *gin.Context) {
	enc, _ := json.Marshal(controller.service.GetAccounts())
	fmt.Fprint(ctx.Writer, string(enc))
}

//GetAccount returns the account identified by account/:id
func (controller *AccountController) GetAccount(ctx *gin.Context) {
	ID, _ := strconv.Atoi(ctx.Param("id"))
	enc, _ := json.Marshal(controller.service.GetAccount(int64(ID)))
	fmt.Fprint(ctx.Writer, string(enc))
}

//UpdateAccount updates the value of the account identified by accounts/:id with the form-value "value" as diffenrence
func (controller *AccountController) UpdateAccount(ctx *gin.Context) {
	ID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		fmt.Fprint(ctx.Writer, "id is NaN")
	}
	value, err := strconv.Atoi(ctx.PostForm("value"))
	if err != nil {
		fmt.Fprint(ctx.Writer, "value is NaN")
	}
	acc, _ := controller.service.UpdateAccount(int64(ID), value)
	enc, _ := json.Marshal(acc)
	fmt.Fprint(ctx.Writer, string(enc))
}
