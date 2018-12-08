package accounts

import (
	"fmt"
	"strconv"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/killingspark/hadibar/src/authStuff"
	"github.com/killingspark/hadibar/src/permissions"
	"github.com/killingspark/hadibar/src/restapi"
)

//AccountController is the controller for accounts
type AccountController struct {
	service *AccountService
}

//NewAccountController creates a new AccountController and initializes the service
func NewAccountController(perms *permissions.Permissions, datadir string) (*AccountController, error) {
	acC := &AccountController{}
	var err error
	acC.service, err = NewAccountService(datadir, perms)
	if err != nil {
		return nil, err
	}
	return acC, nil
}

//GetAccounts gets all existing accounts, that the logged in account is allowed to see
func (controller *AccountController) GetAccounts(ctx *gin.Context) {
	info, err := authStuff.GetLoginInfoFromCtx(ctx)
	if err != nil {
		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}

	accs, err := controller.service.GetAccounts(info.Name)
	if err != nil {
		log.WithFields(log.Fields{"user": info.Name}).WithError(err).Error("Account Error GetAll")

		response, _ := restapi.NewErrorResponse("Could not get accounts because: " + err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}
	response, _ := restapi.NewOkResponse(accs).Marshal()
	fmt.Fprint(ctx.Writer, string(response))
	ctx.Next()

}

//GetAccount returns the account identified by the ID in the query
func (controller *AccountController) GetAccount(ctx *gin.Context) {
	ID, ok := ctx.GetQuery("id")
	if !ok {
		log.WithFields(log.Fields{"URL": ctx.Request.URL.String()}).Warn("No ID found in query")

		response, _ := restapi.NewErrorResponse("no id given in query").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}

	info, err := authStuff.GetLoginInfoFromCtx(ctx)
	if err != nil {
		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}

	acc, err := controller.service.GetAccount(ID, info.Name)
	if err != nil {
		log.WithFields(log.Fields{"user": info.Name}).WithError(err).Error("Account Error Get")

		response, _ := restapi.NewErrorResponse("Error getting account: " + err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}

	response, _ := restapi.NewOkResponse(acc).Marshal()
	fmt.Fprint(ctx.Writer, string(response))
	ctx.Next()
}

//UpdateAccount updates the the account identified by the ID int the query
func (controller *AccountController) UpdateAccount(ctx *gin.Context) {
	ID, ok := ctx.GetQuery("id")
	if !ok {
		log.WithFields(log.Fields{"URL": ctx.Request.URL.String()}).Warn("No ID found in query")

		errResp, _ := restapi.NewErrorResponse("No ID given").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	value, err := strconv.Atoi(ctx.PostForm("value"))
	if err != nil {
		log.WithFields(log.Fields{"URL": ctx.Request.URL.String()}).Warn("No int value found in postform")

		response, _ := restapi.NewErrorResponse("No valid diff value").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}

	info, err := authStuff.GetLoginInfoFromCtx(ctx)
	if err != nil {
		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}

	acc, err := controller.service.UpdateAccount(ID, info.Name, value)
	if err != nil {
		log.WithFields(log.Fields{"user": info.Name}).WithError(err).Error("Account Error Update")

		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}
	response, _ := restapi.NewOkResponse(acc).Marshal()
	fmt.Fprint(ctx.Writer, string(response))
	ctx.Next()
}

//GivePermissionToUser allows an other user to see/alter this account
func (controller *AccountController) GivePermissionToUser(ctx *gin.Context) {
	ID, ok := ctx.GetQuery("id")
	if !ok {
		log.WithFields(log.Fields{"URL": ctx.Request.URL.String()}).Warn("No ID found in query")

		errResp, _ := restapi.NewErrorResponse("No ID given").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	newowner := ctx.PostForm("newowner")
	if newowner == "" {
		log.WithFields(log.Fields{"URL": ctx.Request.URL.String()}).Warn("No newowner found in postform")

		response, _ := restapi.NewErrorResponse("No newowner in postform").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}

	info, err := authStuff.GetLoginInfoFromCtx(ctx)
	if err != nil {
		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}

	err = controller.service.GivePermissionToUser(ID, info.Name, newowner, permissions.CRUD)
	if err != nil {
		log.WithFields(log.Fields{"user": info.Name}).WithError(err).Error("Account Error GivePermission")

		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}
	response, _ := restapi.NewOkResponse("").Marshal()
	fmt.Fprint(ctx.Writer, string(response))
	ctx.Next()
}

//NewAccount creates a new account and chooses a new ID
func (controller *AccountController) NewAccount(ctx *gin.Context) {
	name, ok := ctx.GetPostForm("name")
	if !ok {
		log.WithFields(log.Fields{"URL": ctx.Request.URL.String()}).Warn("No name found in postform")

		errResp, _ := restapi.NewErrorResponse("No name given").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	info, err := authStuff.GetLoginInfoFromCtx(ctx)
	if err != nil {
		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}

	acc, err := controller.service.CreateAdd(name, info.Name, permissions.CRUD)
	if err != nil {
		log.WithFields(log.Fields{"user": info.Name}).WithError(err).Error("Account Error New")

		response, _ := restapi.NewErrorResponse("Couldn't create account: " + err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}
	response, _ := restapi.NewOkResponse(acc).Marshal()
	fmt.Fprint(ctx.Writer, string(response))
	ctx.Next()
}

//DeleteAccount deletes the account identified by the ID in the query
func (controller *AccountController) DeleteAccount(ctx *gin.Context) {
	ID, ok := ctx.GetQuery("id")
	if !ok {
		log.WithFields(log.Fields{"URL": ctx.Request.URL.String()}).Warn("No ID found in query")

		errResp, _ := restapi.NewErrorResponse("No ID given").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	info, err := authStuff.GetLoginInfoFromCtx(ctx)
	if err != nil {
		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}

	if err := controller.service.DeleteAccount(ID, info.Name); err == nil {
		response, _ := restapi.NewOkResponse("").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Next()
	} else {
		log.WithFields(log.Fields{"user": info.Name}).WithError(err).Error("Account Error Delete")

		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}
}

//DoTransaction executes a transaction between the accounts identified by their IDs in the query
func (controller *AccountController) DoTransaction(ctx *gin.Context) {
	sourceID, ok := ctx.GetPostForm("sourceid")
	if !ok {
		log.WithFields(log.Fields{"URL": ctx.Request.URL.String()}).Warn("No SourceID found in postform")

		errResp, _ := restapi.NewErrorResponse("No sourceID given").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	targetID, ok := ctx.GetPostForm("targetid")
	if !ok {
		log.WithFields(log.Fields{"URL": ctx.Request.URL.String()}).Warn("No TargetID found in postform")

		errResp, _ := restapi.NewErrorResponse("No targetID given").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	amount, err := strconv.Atoi(ctx.PostForm("amount"))
	if err != nil {
		log.WithFields(log.Fields{"URL": ctx.Request.URL.String()}).Warn("No int amount found in postform")

		response, _ := restapi.NewErrorResponse("No valid diff value").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}

	info, err := authStuff.GetLoginInfoFromCtx(ctx)
	if err != nil {
		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}

	if err := controller.service.Transaction(sourceID, targetID, info.Name, amount); err == nil {
		response, _ := restapi.NewOkResponse("").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Next()
	} else {
		log.WithFields(log.Fields{"user": info.Name}).WithError(err).Error("Transaction Error")

		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}
}
