package controllers

import (
	"net/http"

	"fmt"

	"github.com/julienschmidt/httprouter"
	"github.com/killingspark/HaDiBar/services"
)

//LoginController is the controller for the logins
type LoginController struct {
	service services.LoginService
}

//MakeLoginController creates a new LoginController and initializes the service
func MakeLoginController() LoginController {
	return LoginController{service: services.LoginService{}}
}

//CheckIdentity checks if the token is valid and then executes the given handle
func (controller *LoginController) CheckIdentity(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var token string
		formToken := r.FormValue("token")
		queryToken := r.URL.Query().Get("token")

		if len(formToken) > 0 && len(queryToken) > 0 {
			if formToken == queryToken {
				token = formToken
			}
		}

		if len(formToken) > 0 && len(queryToken) <= 0 {
			token = formToken
		}

		if len(queryToken) > 0 && len(formToken) <= 0 {
			token = queryToken
		}

		if controller.service.IsTokenValid(token) {
			h(w, r, ps)
		} else {
			// Request Basic Authentication otherwise
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

//NewTokenWithCredentials returns a new token if the credentials (in the formvalues) "name" and "password" are valid
func (controller *LoginController) NewTokenWithCredentials(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	name := r.FormValue("name")
	password := r.FormValue("password")

	var tk, ok = controller.service.RequestToken(name, password)
	if !ok {
		fmt.Fprint(w, "NOPE")
	}

	fmt.Fprint(w, "\""+tk+"\"")
}
