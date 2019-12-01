package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/killingspark/hadibar/pkg/accounts"
	"github.com/killingspark/hadibar/pkg/admin"
	"github.com/killingspark/hadibar/pkg/authStuff"
	"github.com/killingspark/hadibar/pkg/beverages"
	"github.com/killingspark/hadibar/pkg/logger"
	"github.com/killingspark/hadibar/pkg/permissions"
	"github.com/killingspark/hadibar/pkg/reports"

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

	viper.SetEnvPrefix("hadibar")
	viper.AutomaticEnv()

	viper.ReadInConfig()

	
	if err := logger.PrepareLoggerFromViper() 
	err != nil {
		panic("Could not setup logger: " + err.Error())
	}
	

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

	dataDir := viper.GetString("DataDir")
	if stats, err := os.Stat(dataDir); err != nil {
		if os.IsNotExist(err) {
			panic("Datadir does not exist")
		}
	} else {
		if !stats.IsDir() {
			panic("Datadir is no directory")
		}
	}

	socketPath := viper.GetString("SocketPath")
	if socketPath != "" {
		if stats, err := os.Stat(socketPath); err != nil {
			if os.IsNotExist(err) {
				panic("Socketpath does not exist")
			}
		} else {
			if !stats.IsDir() {
				panic("Socketpath is no directory")
			}
		}
	}

	adminTCPAddr := viper.GetString("AdminTcpAddr")
	adminTLSCertPath := viper.GetString("TlsCertPath")
	adminTLSKeyPath := viper.GetString("TlsKeyPath")
	adminTLSCaCertPath := viper.GetString("TlsCaCertPath")
	adminTLSClientCertReq := viper.GetBool("TlsRequireClientCert")

	if adminTCPAddr != "" && socketPath != "" {
		panic("Only one adminserver should be started")
	}

	//////
	// USERS
	//////
	usrRepo, err := authStuff.NewUserRepo(dataDir)
	if err != nil {
		panic(err.Error())
	}

	auth := authStuff.NewAuth(usrRepo, viper.GetInt("SessionTTL"))
	if err != nil {
		panic(err.Error())
	}

	perms, err := permissions.NewPermissions(dataDir)
	if err != nil {
		panic(err.Error())
	}

	//////
	// BEVERAGES
	//////
	br, err := beverages.NewBeverageRepo(dataDir)
	if err != nil {
		panic(err.Error())
	}
	bs := beverages.NewBeverageService(br, perms)
	bc := beverages.NewBeverageController(bs)

	//////
	// ACCOUNTS
	//////
	acr, err := accounts.NewAccountRepo(dataDir)
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

	if socketPath != "" {
		adminSocketPath := socketPath + "/admin.socket"
		os.Remove(adminSocketPath)
		admnSrvr, err := admin.NewUnixAdminServer(adminSocketPath, usrRepo, acr, br, perms)
		if err != nil {
			panic(err.Error())
		}
		go admnSrvr.StartAccepting()
	}
	if adminTCPAddr != "" {
		if adminTLSKeyPath != "" {
			admnSrvr, err := admin.NewTlsAdminServer(adminTCPAddr, adminTLSCertPath, adminTLSKeyPath, adminTLSCaCertPath, adminTLSClientCertReq, usrRepo, acr, br, perms)
			if err != nil {
				panic(err.Error())
			}
			go admnSrvr.StartAccepting()
		} else {
			admnSrvr, err := admin.NewTcpAdminServer(adminTCPAddr, usrRepo, acr, br, perms)
			if err != nil {
				panic(err.Error())
			}
			go admnSrvr.StartAccepting()
		}

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
	makeUserUpdateRoutes(floorSpecificGroup, lc)

	portStr := ":" + viper.GetString("Port")
	if viper.GetInt("Port") <= 0 {
		panic("Port is not a valid port number")
	}
	log.Fatal(http.ListenAndServe(portStr, router))
}
