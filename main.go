package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/killingspark/HaDiBar/controllers"
)

func makeBeverageRoutes(router *httprouter.Router, lc *controllers.LoginController, bc *controllers.BeverageController) {
	router.GET("/beverages", lc.CheckIdentity(bc.GetBeverages))
	router.GET("/beverage/:id", lc.CheckIdentity(bc.GetBeverage))
	router.POST("/beverage/:id", lc.CheckIdentity(bc.UpdateBeverage))
	router.DELETE("/beverage/:id", lc.CheckIdentity(bc.DeleteBeverage))
	router.PUT("/newbeverage", lc.CheckIdentity(bc.NewBeverage))
}

func makeAccountRoutes(router *httprouter.Router, lc *controllers.LoginController, ac *controllers.AccountController) {
	router.GET("/accounts", lc.CheckIdentity(ac.GetAccounts))
	router.GET("/account/:id", lc.CheckIdentity(ac.GetAccount))
	router.POST("/account/:id", lc.CheckIdentity(ac.UpdateAccount))
}

func makeLoginRoutes(router *httprouter.Router, lc *controllers.LoginController) {
	router.GET("/login/token", lc.NewTokenWithCredentials)
}

func main() {
	router := httprouter.New()

	//app ist unter /app erreichbar und served das build verzeichnis von react
	router.ServeFiles("/app/*filepath", http.Dir("app"))
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.Redirect(w, r, "/app", 300)
	})

	bc := controllers.MakeBeverageController()
	ac := controllers.MakeAccountController()
	lc := controllers.MakeLoginController()

	makeBeverageRoutes(router, &lc, &bc)
	makeAccountRoutes(router, &lc, &ac)
	makeLoginRoutes(router, &lc)
	http.ListenAndServe(":8080", router)
}
