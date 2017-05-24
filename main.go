package main

import (
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/killingspark/HaDiBar/controllers"
	"github.com/killingspark/HaDiBar/services"
)

func makeBeverageRoutes(router *gin.Engine, bc *controllers.BeverageController) {
	bevGroup := router.Group("/beverage")
	bevGroup.GET("/:id", bc.GetBeverage)
	bevGroup.POST("/:id", bc.UpdateBeverage)
	bevGroup.DELETE("/:id", bc.DeleteBeverage)
	bevGroup.PUT("/new", bc.NewBeverage)
	bevGroup.GET("/", bc.GetBeverages)
}

func makeAccountRoutes(router *gin.Engine, ac *controllers.AccountController) {
	accGroup := router.Group("/account")
	accGroup.GET("/", ac.GetAccounts)
	accGroup.GET("/:id", ac.GetAccount)
	accGroup.POST("/:id", ac.UpdateAccount)
}

func makeLoginRoutes(router *gin.Engine, lc *controllers.LoginController) {
	router.GET("/login", lc.NewTokenWithCredentials)
	router.GET("/logout", lc.LogOut)
	//used to get an initial session id if wished
	router.GET("/session", func(c *gin.Context) {})
}

//CheckSession checks if the token is valid and then executes the given handle
func CheckSession(ss *services.SessionService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID := ctx.Request.Header.Get("sessionID")

		if sessionID == "" {
			println("no session header found. Adding new one")
			ctx.Writer.Header().Set("sessionID", ss.MakeSessionID())
		} else {
			ctx.Writer.Header().Set("sessionID", sessionID)
			println("call from session: " + sessionID)
		}
		ctx.Writer.WriteHeader(http.StatusCreated)
		ctx.Next()
	}
}

func main() {
	router := gin.New()

	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(300, "/app")
	})
	router.StaticFS("/app", http.Dir("app"))

	ss := services.MakeSessionService()
	bc := controllers.MakeBeverageController()
	ac := controllers.MakeAccountController()
	lc := controllers.MakeLoginController(&ss)

	router.Use(CheckSession(&ss))

	makeBeverageRoutes(router, &bc)
	makeAccountRoutes(router, &ac)
	makeLoginRoutes(router, &lc)

	log.Fatal(http.ListenAndServe(":8080", router))
}
