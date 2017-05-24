package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/killingspark/HaDiBar/services"
)

//LoginController is the controller for the logins
type LoginController struct {
	loginservice   services.LoginService
	sessionservice *services.SessionService
}

//MakeLoginController creates a new LoginController and initializes the service
func MakeLoginController(aSs *services.SessionService) LoginController {
	return LoginController{loginservice: services.LoginService{}, sessionservice: aSs}
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

//NewTokenWithCredentials returns a new token if the credentials (in the formvalues) "name" and "password" are valid
func (controller *LoginController) LogOut(ctx *gin.Context) {
	sessionID := ctx.Request.Header.Get("sessionID")
	session, err := controller.sessionservice.GetSession(sessionID)
	if err != nil {
		return
	}

	session.Token = ""
	fmt.Fprint(ctx.Writer, "OK")
}
