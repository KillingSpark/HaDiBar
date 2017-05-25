package accounts

import (
	"fmt"

	"github.com/gin-gonic/gin"
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

//NewTokenWithCredentials returns a new token if the credentials (in the formvalues) "name" and "password" are valid
func (controller *LoginController) NewTokenWithCredentials(ctx *gin.Context) {
	name := ctx.PostForm("name")
	password := ctx.PostForm("password")

	var tk, ok = controller.loginservice.RequestToken(name, password)
	if !ok {
		fmt.Fprint(ctx.Writer, "NOPE")
	} else {
		sessionID := ctx.Request.Header.Get("sessionID")
		session, err := controller.sessionservice.GetSession(sessionID)
		if err != nil {
			return
		}

		session.Token = tk
		fmt.Fprint(ctx.Writer, "OK")
	}
}

//LogOut uncouples the usersession from a token
func (controller *LoginController) LogOut(ctx *gin.Context) {
	sessionID := ctx.Request.Header.Get("sessionID")
	session, err := controller.sessionservice.GetSession(sessionID)
	if err != nil {
		return
	}

	session.Token = ""
	fmt.Fprint(ctx.Writer, "OK")
}
