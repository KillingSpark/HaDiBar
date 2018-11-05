package main

import (
	"net/http"
	"os"

	"log"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/killingspark/HaDiBar/accounts"
	"github.com/killingspark/HaDiBar/authStuff"
	"github.com/killingspark/HaDiBar/beverages"
	"github.com/killingspark/HaDiBar/logger"
	"github.com/killingspark/HaDiBar/settings"
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
	accGroup.POST("/new", ac.NewAccount)
}

func makeLoginRoutes(router *gin.RouterGroup, lc *authStuff.LoginController) {
	router.POST("/session/login", lc.Login)
	router.POST("/session/logout", lc.LogOut)
	//used to get an initial session id if wished
	router.GET("/session/getid", lc.NewSession)
}

func main() {
	settings.ReadSettings()
	if len(os.Args) == 1 {
		startServer()
	} else {
		if os.Args[1] == "addaccount" {
			addAccount()
		}
	}
}

func addAccount() {
	if len(os.Args) != 6 {
		println("Wrong args. Use: accID name group value")
		return
	}
	acc := &accounts.Account{}

	var err error
	acc.ID = os.Args[2]
	acc.Owner = accounts.AccountOwner{Name: os.Args[3]}
	acc.Groups = []*accounts.AccountGroup{&accounts.AccountGroup{GroupID: os.Args[4]}}
	acc.Value, err = strconv.Atoi(os.Args[5])
	if err != nil {
		print("Cant parse value")
		return
	}

	acs, err := accounts.NewAccountService(settings.S.DataDir)
	if err != nil {
		println(err.Error())
		return
	}
	err = acs.Add(acc)
	if err != nil {
		println(err.Error())
		return
	}
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
		panic(err)
	}

	bc, err := beverages.NewBeverageController()
	if err != nil {
		panic(err)
	}
	ac, err := accounts.NewAccountController()
	if err != nil {
		panic(err)
	}
	lc := authStuff.NewLoginController(auth)

	//router.Use(sessMan.CheckSession)
	apiGroup := router.Group("/api")
	floorSpecificGroup := apiGroup.Group("/f")
	floorSpecificGroup.Use(auth.CheckSession)
	floorSpecificGroup.Use(auth.CheckLoginStatus)

	makeBeverageRoutes(floorSpecificGroup, bc)
	makeAccountRoutes(floorSpecificGroup, ac)
	makeLoginRoutes(apiGroup, lc)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(settings.S.Port), router))
}
