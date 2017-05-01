package controllers

import (
	"encoding/json"
	"net/http"

	"fmt"

	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/killingspark/beverages/services"
)

//BeverageController : Controller for the Beverages
type BeverageController struct {
	service services.IBeverageService
}

func MakeBeverageController() BeverageController {
	var bc BeverageController
	s := services.MakeSQLiteBeverageService()
	bc.service = &s //needed this indirection because the Methods are defined for pointers
	return bc
}

func (this *BeverageController) GetBeverages(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	enc, _ := json.Marshal(this.service.GetBeverages())
	fmt.Fprint(w, string(enc))
}

func (this *BeverageController) GetBeverage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ID, err := strconv.Atoi(ps.ByName("id"))

	if err != nil {
		fmt.Fprint(w, "NOPE")
	}

	bev, ok := this.service.GetBeverage(int64(ID))
	if ok {
		enc, _ := json.Marshal(bev)
		fmt.Fprint(w, string(enc))
	} else {
		fmt.Fprint(w, "NOPE")
	}
}

func (this *BeverageController) NewBeverage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	nv, _ := strconv.Atoi(r.FormValue("value"))
	nb, _ := this.service.NewBeverage(r.FormValue("name"), nv)
	enc, _ := json.Marshal(nb)

	fmt.Fprint(w, string(enc))
}

func (this *BeverageController) UpdateBeverage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ID, err := strconv.Atoi(ps.ByName("id"))

	if err != nil {
		fmt.Fprint(w, "NOPE")
	}
	if err != nil {
		fmt.Fprint(w, "NOPE")
		return
	}
	nv, _ := strconv.Atoi(r.FormValue("value"))
	nn := r.FormValue("name")
	nb, _ := this.service.UpdateBeverage(int64(ID), nn, nv)
	enc, _ := json.Marshal(nb)
	fmt.Fprint(w, string(enc))
}

func (this *BeverageController) DeleteBeverage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ID, err := strconv.Atoi(ps.ByName("id"))

	if err != nil {
		fmt.Fprint(w, "NOPE")
	}
	if err != nil {
		fmt.Fprint(w, "NOPE")
		return
	}
	this.service.DeleteBeverage(int64(ID))

}
