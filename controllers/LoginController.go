package controllers

import (
	"net/http"

	"fmt"

	"github.com/julienschmidt/httprouter"
	"github.com/killingspark/HaDiBar/services"
)

//LoginController is the controller for the logins
type LoginController struct {
	loginservice   services.LoginService
	sessionservice services.SessionService
}

//MakeLoginController creates a new LoginController and initializes the service
func MakeLoginController() LoginController {
	return LoginController{loginservice: services.LoginService{}, sessionservice: services.MakeSessionService()}
}

//CheckIdentity checks if the token is valid and then executes the given handle
func (controller *LoginController) CheckIdentity(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		sessionID := r.Header.Get("sessionID")

		if sessionID == "" {
			println("no session header found. Adding new one")
			w.Header().Set("sessionID", controller.sessionservice.MakeSessionID())
		} else {
			w.Header().Set("sessionID", sessionID)
			println("call from session: " + sessionID)
		}
		w.WriteHeader(http.StatusCreated)
		h(w, r, ps)
	}
}

//NewTokenWithCredentials returns a new token if the credentials (in the formvalues) "name" and "password" are valid
func (controller *LoginController) NewTokenWithCredentials(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	name := r.FormValue("name")
	password := r.FormValue("password")

	var tk, ok = controller.loginservice.RequestToken(name, password)
	if !ok {
		fmt.Fprint(w, "NOPE")
	} else {
		sessionCookie, err := r.Cookie("sessionID")
		if err != nil {
			return
		}

		sessionID := sessionCookie.Value
		session, err := controller.sessionservice.GetSession(sessionID)
		if err != nil {
			return
		}

		session.Token = tk
		fmt.Fprint(w, "OK")

	}
}
