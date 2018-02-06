package accounts

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/killingspark/HaDiBar/logger"
	"github.com/killingspark/HaDiBar/restapi"
	"github.com/killingspark/HaDiBar/sessions"
)

//LoginController is the controller for the logins
type LoginController struct {
	loginservice *LoginService
	sesMan       *sessions.SessionManager
}

//NewLoginController creates a new LoginController and initializes the service
func NewLoginController(aSs *sessions.SessionManager) *LoginController {
	return &LoginController{loginservice: NewLoginService(), sesMan: aSs}
}

//Login returns a new token if the credentials (in the formvalues) "name" and "password" are valid
func (controller *LoginController) Login(ctx *gin.Context) {
	name := ctx.PostForm("name")
	password := ctx.PostForm("password")
	sessionID := ctx.Request.Header.Get("sessionID")

	logger.Logger.Debug("Requesting token for: " + name)
	var tk, ok = controller.loginservice.RequestToken(name, password)
	entity, _ := controller.loginservice.GetEntityFromToken(tk)
	floor := entity.Floor

	if !ok {
		response, _ := restapi.NewErrorResponse("credentials rejected").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		logger.Logger.Debug(sessionID + " faild to log in as: " + name)
	} else {
		controller.sesMan.SetSessionToken(sessionID, tk)
		controller.sesMan.SetSessionName(sessionID, name)
		controller.sesMan.SetSessionFloor(sessionID, floor)
		response, _ := restapi.NewOkResponse("").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		logger.Logger.Debug(sessionID + " logged in as: " + name)
	}
}

//LogOut uncouples the usersession from a token
func (controller *LoginController) LogOut(ctx *gin.Context) {
	sessionID := ctx.Request.Header.Get("sessionID")
	controller.sesMan.SetSessionToken(sessionID, "")
	controller.sesMan.SetSessionName(sessionID, "")
	controller.sesMan.SetSessionFloor(sessionID, "")
	response, _ := restapi.NewOkResponse("").Marshal()
	fmt.Fprint(ctx.Writer, string(response))

	logger.Logger.Debug(sessionID + " logged out")
}
