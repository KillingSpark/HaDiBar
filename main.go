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

	"golang.org/x/crypto/sha3"
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
		if os.Args[1] == "adduser" {
			addUser()
		}
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
	acc.ID, err = strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		print("Cant parse accountID")
		return
	}
	acc.Owner = accounts.AccountOwner{Name: os.Args[3]}
	acc.Group = accounts.AccountGroup{GroupID: os.Args[4]}
	acc.Value, err = strconv.Atoi(os.Args[5])
	if err != nil {
		print("Cant parse value")
		return
	}

	acs := accounts.NewAccountService()
	acs.Load()
	err = acs.Add(acc)
	if err != nil {
		println(err.Error())
		return
	}
	acs.Save()
}

func addUser() {
	if len(os.Args) != 6 {
		println("Wrong args. Use: name group salt passwd")
		return
	}
	info := &authStuff.LoginInfo{}
	info.Name = os.Args[2]
	info.GroupID = os.Args[3]
	info.Salt = os.Args[4]
	info.Pwhash = authStuff.SaltPw(sha3.New256(), os.Args[5], info.Salt)

	lh := authStuff.NewJsonUserDatabase()
	lh.Load()
	err := lh.Add(info)
	if err != nil {
		print(err.Error())
	} else {
		lh.Save()
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

	auth := authStuff.NewAuth()

	bc := beverages.NewBeverageController()
	ac := accounts.NewAccountController(auth)
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
