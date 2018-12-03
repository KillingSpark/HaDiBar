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

//NewSession creates a new session id and writes it to as an answer
func (controller *LoginController) NewSession(ctx *gin.Context) {
	id := controller.auth.AddNewSession()
	fmt.Fprint(ctx.Writer, id)
	ctx.Next()
}

//SetEmail sets the email for a logged in user
func (controller *LoginController) SetEmail(ctx *gin.Context) {
	email := ctx.PostForm("email")
	if email == "" {
		response, _ := restapi.NewErrorResponse("No email given").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}

	info, err := GetLoginInfoFromCtx(ctx)
	if err != nil {
		response, _ := restapi.NewErrorResponse("No Logininfo found. This is an internl error that should never happen").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}
	controller.auth.ls.SetEmail(info.Name, email)
	response, _ := restapi.NewOkResponse("").Marshal()
	fmt.Fprint(ctx.Writer, string(response))
	ctx.Next()
}

//GetUser sends the info for a logged in user
func (controller *LoginController) GetUser(ctx *gin.Context) {
	info, err := GetLoginInfoFromCtx(ctx)
	if err != nil {
		response, _ := restapi.NewErrorResponse("No Logininfo found. This is an internl error that should never happen").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}

	type userresp struct {
		Email string
		Name  string
	}
	response, _ := restapi.NewOkResponse(userresp{Name: info.Name, Email: info.Email}).Marshal()
	fmt.Fprint(ctx.Writer, string(response))
	ctx.Next()
}

//Login checks whether "name" and "password" are valid and updates the logininfo if so
func (controller *LoginController) Login(ctx *gin.Context) {
	name := ctx.PostForm("name")
	password := ctx.PostForm("password")
	sessionID := ctx.Request.Header.Get("sessionID")

	err := controller.auth.LogIn(sessionID, name, password)

	if err != nil {
		response, _ := restapi.NewErrorResponse("credentials rejected: " + err.Error()).Marshal()
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

//LogOut uncouples the session from the logininfo
func (controller *LoginController) LogOut(ctx *gin.Context) {
	sessionID := ctx.Request.Header.Get("sessionID")
	controller.auth.LogOut(sessionID)
	response, _ := restapi.NewOkResponse("").Marshal()
	fmt.Fprint(ctx.Writer, string(response))

	logger.Logger.Debug(sessionID + " logged out")
	ctx.Next()
}
