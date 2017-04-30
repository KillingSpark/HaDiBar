package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/killingspark/beverages/services"
)

type AccountController struct {
	service services.AccountService
}

func MakeAccountController() AccountController {
	var acC AccountController
	acC.service = services.MakeAccountService()
	return acC
}

func (this *AccountController) GetAccounts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	enc, _ := json.Marshal(this.service.GetAccounts())
	fmt.Fprint(w, string(enc))
}

func (this *AccountController) GetAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ID, _ := strconv.Atoi(ps.ByName("id"))
	enc, _ := json.Marshal(this.service.GetAccount(int64(ID)))
	fmt.Fprint(w, string(enc))
}

func (this *AccountController) UpdateAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		fmt.Fprintln(w, "id is NaN")
	}
	value, err := strconv.Atoi(r.FormValue("value"))
	if err != nil {
		fmt.Fprintln(w, "value is NaN")
	}
	acc, _ := this.service.UpdateAccount(int64(ID), value)
	enc, _ := json.Marshal(acc)
	fmt.Fprint(w, string(enc))
}
