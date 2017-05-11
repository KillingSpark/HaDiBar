package controllers

import (
	"encoding/json"
	"net/http"

	"fmt"

	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/killingspark/HaDiBar/services"
)

//BeverageController : Controller for the Beverages
type BeverageController struct {
	service *services.SQLiteBeverageService
}

//MakeBeverageController creates a new BeverageController and initializes its service
func MakeBeverageController() BeverageController {
	var bc BeverageController
	s := services.MakeSQLiteBeverageService()
	bc.service = &s //needed controller indirection because the Methods are defined for pointers
	return bc
}

//GetBeverages responds with all existing Beverages
func (controller *BeverageController) GetBeverages(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	enc, _ := json.Marshal(controller.service.GetBeverages())
	fmt.Fprint(w, string(enc))
}

//GetBeverage responds with the beverage identified by beverage/:id
func (controller *BeverageController) GetBeverage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ID := ps.ByName("id")

	bev, ok := controller.service.GetBeverage(ID)
	if ok {
		enc, _ := json.Marshal(bev)
		fmt.Fprint(w, string(enc))
	} else {
		fmt.Fprint(w, "NOPE")
	}
}

//NewBeverage creates a new beverage with the given form-values "value" and "name" and returns it
func (controller *BeverageController) NewBeverage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	nv, _ := strconv.Atoi(r.FormValue("value"))
	nb, _ := controller.service.NewBeverage(r.FormValue("name"), nv)
	enc, _ := json.Marshal(nb)

	fmt.Fprint(w, string(enc))
}

//UpdateBeverage updates the beverage identified by /beverage/:id with the given form-values "value" and "name" and returns it
func (controller *BeverageController) UpdateBeverage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ID := ps.ByName("id")

	nv, _ := strconv.Atoi(r.FormValue("value"))
	nn := r.FormValue("name")
	nb, _ := controller.service.UpdateBeverage(ID, nn, nv)
	enc, _ := json.Marshal(nb)
	fmt.Fprint(w, string(enc))
}

//DeleteBeverage deletes the beverage identified by /beverage/:id and responds with a YEP/NOPE
func (controller *BeverageController) DeleteBeverage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ID := ps.ByName("id")

	if controller.service.DeleteBeverage(ID) {
		fmt.Fprint(w, "YEP")
	} else {
		fmt.Fprint(w, "NOPE")
	}

}
