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
	loginservice   *LoginService
	sessionservice *sessions.SessionManager
}

//NewLoginController creates a new LoginController and initializes the service
func NewLoginController(aSs *sessions.SessionManager) *LoginController {
	return &LoginController{loginservice: &LoginService{}, sessionservice: aSs}
}

//Login returns a new token if the credentials (in the formvalues) "name" and "password" are valid
func (controller *LoginController) Login(ctx *gin.Context) {
	name := ctx.PostForm("name")
	password := ctx.PostForm("password")

	logger.Logger.Debug("Requesting token for: " + name)
	var tk, ok = controller.loginservice.RequestToken(name, password)
	logger.Logger.Debug("Received token for: " + name + " : " + tk)
	if !ok {
		response, _ := restapi.NewErrorResponse("credentials rejected").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
	} else {
		sessionID := ctx.Request.Header.Get("sessionID")
		controller.sessionservice.SetSessionToken(sessionID, tk)
		controller.sessionservice.SetSessionName(sessionID, name)
		response, _ := restapi.NewOkResponse("").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
	}
}

//LogOut uncouples the usersession from a token
func (controller *LoginController) LogOut(ctx *gin.Context) {
	sessionID := ctx.Request.Header.Get("sessionID")

	controller.sessionservice.SetSessionToken(sessionID, "")
	controller.sessionservice.SetSessionName(sessionID, "")
	response, _ := restapi.NewOkResponse("").Marshal()
	fmt.Fprint(ctx.Writer, string(response))
}
