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
	return LoginController{}
}

//NewTokenWithCredentials returns a new token if the credentials (in the formvalues) "name" and "password" are valid
func (controller *LoginController) NewTokenWithCredentials(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	name := r.FormValue("name")
	password := r.FormValue("password")

	var tk, ok = controller.service.RequestToken(name, password)
	if !ok {
		fmt.Fprint(w, "NOPE")
	}

	fmt.Fprint(w, tk)
}
