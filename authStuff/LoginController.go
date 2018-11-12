package authStuff

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/killingspark/hadibar/logger"
	"github.com/killingspark/hadibar/restapi"
)

//LoginController is the controller for the logins
type LoginController struct {
	auth *Auth
}

//NewLoginController creates a new LoginController and initializes the service
func NewLoginController(auth *Auth) *LoginController {
	return &LoginController{auth: auth}
}

func (controller *LoginController) NewSession(ctx *gin.Context) {
	id := controller.auth.AddNewSession()
	fmt.Fprint(ctx.Writer, id)
	ctx.Next()
}

//Login returns a new token if the credentials (in the formvalues) "name" and "password" are valid
func (controller *LoginController) Login(ctx *gin.Context) {
	name := ctx.PostForm("name")
	password := ctx.PostForm("password")
	sessionID := ctx.Request.Header.Get("sessionID")

	err := controller.auth.LogIn(sessionID, name, password)

	if err != nil {
		response, _ := restapi.NewErrorResponse("credentials rejected: " + sessionID).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		logger.Logger.Debug(sessionID + " failed to log in as: " + name)
		ctx.Abort()
		return
	}
	response, _ := restapi.NewOkResponse("").Marshal()
	fmt.Fprint(ctx.Writer, string(response))
	logger.Logger.Debug(sessionID + " logged in as: " + name)
	ctx.Next()
}

//LogOut uncouples the usersession from a token
func (controller *LoginController) LogOut(ctx *gin.Context) {
	sessionID := ctx.Request.Header.Get("sessionID")
	controller.auth.LogOut(sessionID)
	response, _ := restapi.NewOkResponse("").Marshal()
	fmt.Fprint(ctx.Writer, string(response))

	logger.Logger.Debug(sessionID + " logged out")
	ctx.Next()
}
