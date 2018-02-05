package accounts

import (
	"encoding/json"
	"fmt"

	"github.com/killingspark/HaDiBar/sessions"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/killingspark/HaDiBar/restapi"
)

//AccountController is the controller for accounts
type AccountController struct {
	service *AccountService
	sesMan  *sessions.SessionManager
}

//NewAccountController creates a new AccountController and initializes the service
func NewAccountController(sm *sessions.SessionManager) *AccountController {
	var acC AccountController
	acC.service = NewAccountService()
	acC.sesMan = sm
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
	sessionID := ctx.Request.Header.Get("sessionID")
	session, err := controller.sesMan.GetSession(sessionID)
	if err != nil {
		enc, _ := json.Marshal(restapi.Response{Status: "ERROR", Response: err.Error()})
		fmt.Fprint(ctx.Writer, string(enc))
	}

	acc, err := controller.service.UpdateAccount(session.Token, int64(ID), value)
	if err != nil {
		enc, _ := json.Marshal(restapi.Response{Status: "ERROR", Response: err.Error()})
		fmt.Fprint(ctx.Writer, string(enc))
	}
	enc, _ := json.Marshal(restapi.Response{Status: "OK", Response: acc})
	fmt.Fprint(ctx.Writer, string(enc))
}
