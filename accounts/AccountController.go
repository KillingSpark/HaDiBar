package accounts

import (
	"fmt"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/killingspark/HaDiBar/authStuff"
	"github.com/killingspark/HaDiBar/restapi"
)

//AccountController is the controller for accounts
type AccountController struct {
	service *AccountService
	auth    *authStuff.Auth
}

//NewAccountController creates a new AccountController and initializes the service
func NewAccountController(auth *authStuff.Auth) *AccountController {
	var acC AccountController
	acC.service = NewAccountService()
	acC.auth = auth
	return &acC
}

//GetAccounts gets all existing accounts
func (controller *AccountController) GetAccounts(ctx *gin.Context) {
	if inter, ok := ctx.Get("logininfo"); ok {
		info, ok := inter.(*authStuff.LoginInfo)
		if ok {
			response, _ := restapi.NewOkResponse(controller.service.GetAccounts(info.GroupID)).Marshal()
			fmt.Fprint(ctx.Writer, string(response))
		} else {
			response, _ := restapi.NewErrorResponse("Something went wrong while processing the username").Marshal()
			fmt.Fprint(ctx.Writer, string(response))
		}
	} else {
		response, _ := restapi.NewErrorResponse("Something went wrong while processing the username").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
	}
}

//GetAccount returns the account identified by account/:id
func (controller *AccountController) GetAccount(ctx *gin.Context) {
	strID, ok := ctx.GetQuery("id")
	if !ok {
		response, _ := restapi.NewErrorResponse("no id in path").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		return
	}
	ID, err := strconv.Atoi(strID)
	if err != nil {
		response, _ := restapi.NewErrorResponse("ID is invalid").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		return
	}

	acc, err := controller.service.GetAccount(int64(ID))
	if err != nil {
		response, _ := restapi.NewErrorResponse("Error getting account: " + err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		return
	}

	response, _ := restapi.NewOkResponse(acc).Marshal()
	fmt.Fprint(ctx.Writer, string(response))
}

//UpdateAccount updates the value of the account identified by accounts/:id with the form-value "value" as diffenrence
func (controller *AccountController) UpdateAccount(ctx *gin.Context) {
	strID, _ := ctx.GetQuery("id")
	ID, err := strconv.Atoi(strID)

	if err != nil {
		response, _ := restapi.NewErrorResponse("No valid account id").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}
	value, err := strconv.Atoi(ctx.PostForm("value"))
	if err != nil {
		response, _ := restapi.NewErrorResponse("No valid diff value").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}
	sessionID := ctx.Request.Header.Get("sessionID")
	info, err := controller.auth.GetSessionInfo(sessionID)
	if err != nil {
		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}

	acc, err := controller.service.UpdateAccount(info, int64(ID), value)
	if err != nil {
		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}
	response, _ := restapi.NewOkResponse(acc).Marshal()
	fmt.Fprint(ctx.Writer, string(response))
	ctx.Next()
}
