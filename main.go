package main

import (
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/killingspark/HaDiBar/accounts"
	"github.com/killingspark/HaDiBar/beverages"
	"github.com/killingspark/HaDiBar/sessions"
)

var (
	sessMan = sessions.NewSessionManager()
)

func makeBeverageRoutes(router *gin.Engine, bc *beverages.BeverageController) {
	bevGroup := router.Group("/beverage")
	bevGroup.GET("/:id", bc.GetBeverage)
	bevGroup.POST("/:id", bc.UpdateBeverage)
	bevGroup.DELETE("/:id", bc.DeleteBeverage)
	bevGroup.PUT("/new", bc.NewBeverage)
	bevGroup.GET("/", bc.GetBeverages)
}

func makeAccountRoutes(router *gin.Engine, ac *accounts.AccountController) {
	accGroup := router.Group("/account")
	accGroup.GET("/", ac.GetAccounts)
	accGroup.GET("/:id", ac.GetAccount)
	accGroup.POST("/:id", ac.UpdateAccount)
}

func makeLoginRoutes(router *gin.Engine, lc *accounts.LoginController) {
	router.GET("/login", lc.NewTokenWithCredentials)
	router.GET("/logout", lc.LogOut)
	//used to get an initial session id if wished
	router.GET("/session", func(c *gin.Context) {})
}

func main() {
	router := gin.New()

	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(300, "/app")
	})
	router.StaticFS("/app", http.Dir("webapp"))

	bc := beverages.NewBeverageController()
	ac := accounts.NewAccountController()
	lc := accounts.NewLoginController(sessMan)

	router.Use(sessMan.CheckSession(sessMan))

	makeBeverageRoutes(router, bc)
	makeAccountRoutes(router, ac)
	makeLoginRoutes(router, lc)

	log.Fatal(http.ListenAndServe(":8080", router))
}
