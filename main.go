package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/killingspark/HaDiBar/controllers"
)

func checkIdentity(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token := r.FormValue("token")

		if token != "" {
			h(w, r, ps)
		} else {
			// Request Basic Authentication otherwise
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func makeBeverageRoutes(router *httprouter.Router, bc controllers.BeverageController) {
	router.GET("/beverages", checkIdentity(bc.GetBeverages))
	router.GET("/beverage/:id", checkIdentity(bc.GetBeverage))
	router.POST("/beverage/:id", checkIdentity(bc.UpdateBeverage))
	router.DELETE("/beverage/:id", checkIdentity(bc.DeleteBeverage))
	router.PUT("/newbeverage", checkIdentity(bc.NewBeverage))
}

func makeAccountRoutes(router *httprouter.Router, ac controllers.AccountController) {
	router.GET("/accounts", checkIdentity(ac.GetAccounts))
	router.GET("/account/:id", checkIdentity(ac.GetAccount))
	router.POST("/account/:id", checkIdentity(ac.UpdateAccount))
}

func makeLoginRoutes(router *httprouter.Router, lc controllers.LoginController) {
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

	makeBeverageRoutes(router, bc)
	makeAccountRoutes(router, ac)
	makeLoginRoutes(router, lc)
	http.ListenAndServe(":8080", router)
}
