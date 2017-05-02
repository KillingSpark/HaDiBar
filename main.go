package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/killingspark/HaDiBar/controllers"
)

func makeBeverageRoutes(router *httprouter.Router, bc controllers.BeverageController) {
	router.GET("/beverages", bc.GetBeverages)
	router.GET("/beverage/:id", bc.GetBeverage)
	router.POST("/beverage/:id", bc.UpdateBeverage)
	router.DELETE("/beverage/:id", bc.DeleteBeverage)
	router.PUT("/newbeverage", bc.NewBeverage)
}

func makeAccountRoutes(router *httprouter.Router, ac controllers.AccountController) {
	router.GET("/accounts", ac.GetAccounts)
	router.GET("/account/:id", ac.GetAccount)
	router.POST("/account/:id", ac.UpdateAccount)
}

func makeLoginRoutes(router *httprouter.Router, lc controllers.LoginController) {
	router.GET("/login/token", lc.NewTokenWithCredentials)
}

func main() {
	router := httprouter.New()
	bc := controllers.MakeBeverageController()
	ac := controllers.MakeAccountController()
	lc := controllers.MakeLoginController()

	makeBeverageRoutes(router, bc)
	makeAccountRoutes(router, ac)
	makeLoginRoutes(router, lc)
	http.ListenAndServe(":8080", router)
}
