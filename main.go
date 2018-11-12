package main

import (
	"net/http"

	"log"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/killingspark/hadibar/accounts"
	"github.com/killingspark/hadibar/authStuff"
	"github.com/killingspark/hadibar/beverages"
	"github.com/killingspark/hadibar/logger"
	"github.com/killingspark/hadibar/permissions"
	"github.com/killingspark/hadibar/reports"
	"github.com/killingspark/hadibar/settings"
)

//making routes seperate for better readability
func makeBeverageRoutes(router *gin.RouterGroup, bc *beverages.BeverageController) {
	bevGroup := router.Group("/beverage")
	bevGroup.GET("/all", bc.GetBeverages)
	bevGroup.GET("/get", bc.GetBeverage)
	bevGroup.POST("/update", bc.UpdateBeverage)
	bevGroup.POST("/addToGroup", bc.GivePermissionToUser)
	bevGroup.PUT("/new", bc.NewBeverage)
	bevGroup.DELETE("/delete", bc.DeleteBeverage)
}

func makeAccountRoutes(router *gin.RouterGroup, ac *accounts.AccountController) {
	accGroup := router.Group("/account")
	accGroup.GET("/all", ac.GetAccounts)
	accGroup.GET("/get", ac.GetAccount)
	accGroup.POST("/update", ac.UpdateAccount)
	accGroup.POST("/addToGroup", ac.GivePermissionToUser)
	accGroup.POST("/transaction", ac.DoTransaction)
	accGroup.PUT("/new", ac.NewAccount)
	accGroup.DELETE("/delete", ac.DeleteAccount)
}

func makeReportRoutes(router *gin.RouterGroup, rc *reports.ReportsController) {
	repGroup := router.Group("/reports")
	repGroup.GET("/accounts", rc.GenerateAccountList)
	repGroup.GET("/beverages", rc.GenerateBeverageMatrix)
}

func makeLoginRoutes(router *gin.RouterGroup, lc *authStuff.LoginController) {
	router.POST("/session/login", lc.Login)
	router.POST("/session/logout", lc.LogOut)
	//used to get an initial session id if wished
	router.GET("/session/getid", lc.NewSession)
}

func main() {
	settings.ReadSettings()
	startServer()
}

func startServer() {
	logger.PrepareLogger()
	router := gin.New()

	//serves the wepapp folder as /app
	router.StaticFS(settings.S.WebappRoute, http.Dir(settings.S.WebappPath))

	//redirect users from / to /app
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(300, settings.S.WebappRoute)
	})

	auth, err := authStuff.NewAuth()
	if err != nil {
		panic(err.Error())
	}

	perms, err := permissions.NewPermissions(settings.S.DataDir)
	if err != nil {
		panic(err.Error())
	}

	bc, err := beverages.NewBeverageController(perms)
	if err != nil {
		panic(err.Error())
	}
	ac, err := accounts.NewAccountController(perms)
	if err != nil {
		panic(err.Error())
	}
	lc := authStuff.NewLoginController(auth)

	rc, err := reports.NewReportsController(perms)
	if err != nil {
		panic(err.Error())
	}

	//router.Use(sessMan.CheckSession)
	apiGroup := router.Group("/api")
	floorSpecificGroup := apiGroup.Group("/f")
	floorSpecificGroup.Use(auth.CheckSession)
	floorSpecificGroup.Use(auth.CheckLoginStatus)

	makeBeverageRoutes(floorSpecificGroup, bc)
	makeAccountRoutes(floorSpecificGroup, ac)
	makeReportRoutes(floorSpecificGroup, rc)
	makeLoginRoutes(apiGroup, lc)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(settings.S.Port), router))
}
