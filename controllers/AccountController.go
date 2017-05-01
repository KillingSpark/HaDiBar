package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/killingspark/beverages/services"
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
func (controller *AccountController) GetAccounts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	enc, _ := json.Marshal(controller.service.GetAccounts())
	fmt.Fprint(w, string(enc))
}

//GetAccount returns the account identified by account/:id
func (controller *AccountController) GetAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ID, _ := strconv.Atoi(ps.ByName("id"))
	enc, _ := json.Marshal(controller.service.GetAccount(int64(ID)))
	fmt.Fprint(w, string(enc))
}

//UpdateAccount updates the value of the account identified by accounts/:id with the form-value "value" as diffenrence
func (controller *AccountController) UpdateAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		fmt.Fprintln(w, "id is NaN")
	}
	value, err := strconv.Atoi(r.FormValue("value"))
	if err != nil {
		fmt.Fprintln(w, "value is NaN")
	}
	acc, _ := controller.service.UpdateAccount(int64(ID), value)
	enc, _ := json.Marshal(acc)
	fmt.Fprint(w, string(enc))
}
