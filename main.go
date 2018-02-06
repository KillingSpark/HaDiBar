package main

import (
	"fmt"
	"net/http"

	"log"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/killingspark/HaDiBar/accounts"
	"github.com/killingspark/HaDiBar/beverages"
	"github.com/killingspark/HaDiBar/logger"
	"github.com/killingspark/HaDiBar/sessions"
	"github.com/killingspark/HaDiBar/settings"
)

var (
	SessMan = sessions.NewSessionManager()
)

//making routes seperate for better readability
func makeBeverageRoutes(router *gin.RouterGroup, bc *beverages.BeverageController) {
	bevGroup := router.Group("/beverage")
	bevGroup.GET("/get", bc.GetBeverage)
	bevGroup.POST("/update", bc.UpdateBeverage)
	bevGroup.DELETE("/delete", bc.DeleteBeverage)
	bevGroup.PUT("/new", bc.NewBeverage)
	bevGroup.GET("/all", bc.GetBeverages)
}

func makeAccountRoutes(router *gin.RouterGroup, ac *accounts.AccountController) {
	accGroup := router.Group("/account")
	accGroup.GET("/all", ac.GetAccounts)
	accGroup.GET("/get", ac.GetAccount)
	accGroup.POST("/update", ac.UpdateAccount)
}

func makeLoginRoutes(router *gin.RouterGroup, lc *accounts.LoginController) {
	router.POST("/session/login", lc.Login)
	router.POST("/session/logout", lc.LogOut)
	//used to get an initial session id if wished
	router.GET("/session/getid", func(ctx *gin.Context) { fmt.Fprint(ctx.Writer, SessMan.MakeSessionID()) })
}

func main() {
	logger.PrepareLogger()
	settings.ReadSettings()
	router := gin.New()

	//serves the wepapp folder as /app
	router.StaticFS(settings.S.WebappRoute, http.Dir(settings.S.WebappPath))

	//redirect users from / to /app
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(300, settings.S.WebappRoute)
	})

	bc := beverages.NewBeverageController()
	ac := accounts.NewAccountController(SessMan)
	lc := accounts.NewLoginController(SessMan)

	//router.Use(sessMan.CheckSession)
	apiGroup := router.Group("/api")
	floorSpecificGroup := apiGroup.Group("/f")
	floorSpecificGroup.Use(SessMan.CheckSession)
	floorSpecificGroup.Use(SessMan.CheckLoginStatus)

	makeBeverageRoutes(floorSpecificGroup, bc)
	makeAccountRoutes(floorSpecificGroup, ac)
	makeLoginRoutes(apiGroup, lc)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(settings.S.Port), router))
}
