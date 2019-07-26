package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/killingspark/hadibar/src/accounts"
	"github.com/killingspark/hadibar/src/admin"
	"github.com/killingspark/hadibar/src/authStuff"
	"github.com/killingspark/hadibar/src/beverages"
	"github.com/killingspark/hadibar/src/logger"
	"github.com/killingspark/hadibar/src/permissions"
	"github.com/killingspark/hadibar/src/reports"

	"github.com/spf13/viper"
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
	repGroup.GET("/transactions", rc.GenerateTransactionList)
}

func makeLoginRoutes(router *gin.RouterGroup, lc *authStuff.LoginController) {
	router.POST("/session/login", lc.Login)
	router.POST("/session/logout", lc.LogOut)
	//used to get an initial session id if wished
	router.GET("/session/getid", lc.NewSession)
}
func makeUserUpdateRoutes(router *gin.RouterGroup, lc *authStuff.LoginController) {
	router.POST("/user/email", lc.SetEmail)
	router.GET("/user/info", lc.GetUser)
}

func main() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/hadibar/")
	viper.AddConfigPath("/etc/hadibar")
	viper.SetConfigName("settings")
	viper.ReadInConfig()

	logger.PrepareLoggerFromViper()

	startServer()
}

func startServer() {
	router := gin.New()

	//serves the wepapp folder as /app
	router.StaticFS(viper.GetString("WebAppRoute"), http.Dir(viper.GetString("WebAppDir")))

	//redirect users from / to /app
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(300, viper.GetString("WebAppRoute"))
	})

	//////
	// USERS
	//////
	usrRepo, err := authStuff.NewUserRepo(viper.GetString("DataDir"))
	if err != nil {
		panic(err.Error())
	}

	auth := authStuff.NewAuth(usrRepo, viper.GetInt("SessionTTL"))
	if err != nil {
		panic(err.Error())
	}

	perms, err := permissions.NewPermissions(viper.GetString("DataDir"))
	if err != nil {
		panic(err.Error())
	}

	//////
	// BEVERAGES
	//////
	br, err := beverages.NewBeverageRepo(viper.GetString("DataDir"))
	if err != nil {
		panic(err.Error())
	}
	bs := beverages.NewBeverageService(br, perms)
	bc := beverages.NewBeverageController(bs)

	//////
	// ACCOUNTS
	//////
	acr, err := accounts.NewAccountRepo(viper.GetString("DataDir"))
	if err != nil {
		panic(err.Error())
	}
	acs := accounts.NewAccountService(acr, perms)
	ac := accounts.NewAccountController(acs)

	lc := authStuff.NewLoginController(auth)

	rc := reports.NewReportsController(bs, acs)
	if err != nil {
		panic(err.Error())
	}

	os.Remove(viper.GetString("SocketPath"))
	admnSrvr, err := admin.NewAdminServer(viper.GetString("SocketPath"), usrRepo, acr, br, perms)
	if err != nil {
		panic(err.Error())
	}
	go admnSrvr.StartAccepting()

	//router.Use(sessMan.CheckSession)
	apiGroup := router.Group("/api")
	floorSpecificGroup := apiGroup.Group("/f")
	floorSpecificGroup.Use(auth.CheckSession)
	floorSpecificGroup.Use(auth.CheckLoginStatus)

	makeBeverageRoutes(floorSpecificGroup, bc)
	makeAccountRoutes(floorSpecificGroup, ac)
	makeReportRoutes(floorSpecificGroup, rc)
	makeLoginRoutes(apiGroup, lc)
	makeUserUpdateRoutes(floorSpecificGroup, lc)

	log.Fatal(http.ListenAndServe(":"+viper.GetString("Port"), router))
}
