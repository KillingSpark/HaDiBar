package beverages

import (
	"encoding/json"

	"fmt"

	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/killingspark/HaDiBar/restapi"
)

//BeverageController : Controller for the Beverages
type BeverageController struct {
	service *SQLiteBeverageService
}

//NewBeverageController creates a new BeverageController and initializes its service
func NewBeverageController() *BeverageController {
	var bc BeverageController
	bc.service = NewSQLiteBeverageService()
	return &bc
}

//GetBeverages responds with all existing Beverages
func (controller *BeverageController) GetBeverages(ctx *gin.Context) {
	enc, _ := json.Marshal(restapi.Response{Status: "OK", Response: controller.service.GetBeverages()})
	fmt.Fprint(ctx.Writer, string(enc))
}

//GetBeverage responds with the beverage identified by beverage/:id
func (controller *BeverageController) GetBeverage(ctx *gin.Context) {
	ID := ctx.Param("id")

	bev, ok := controller.service.GetBeverage(ID)
	if ok {
		enc, _ := json.Marshal(restapi.Response{Status: "OK", Response: bev})
		fmt.Fprint(ctx.Writer, string(enc))
	} else {
		fmt.Fprint(ctx.Writer, "{\"status\":\"ERROR\"}")
	}
}

//NewBeverage creates a new beverage with the given form-values "value" and "name" and returns it
func (controller *BeverageController) NewBeverage(ctx *gin.Context) {
	nv, _ := strconv.Atoi(ctx.PostForm("value"))
	nb, _ := controller.service.NewBeverage(ctx.PostForm("name"), nv)
	enc, _ := json.Marshal(restapi.Response{Status: "OK", Response: nb})

	fmt.Fprint(ctx.Writer, string(enc))
}

//UpdateBeverage updates the beverage identified by /beverage/:id with the given form-values "value" and "name" and returns it
func (controller *BeverageController) UpdateBeverage(ctx *gin.Context) {
	ID := ctx.Param("id")

	nv, _ := strconv.Atoi(ctx.PostForm("value"))
	nn := ctx.PostForm("name")
	nb, _ := controller.service.UpdateBeverage(ID, nn, nv)
	enc, _ := json.Marshal(restapi.Response{Status: "OK", Response: nb})

	fmt.Fprint(ctx.Writer, string(enc))
}

//DeleteBeverage deletes the beverage identified by /beverage/:id and responds with a YEP/NOPE
func (controller *BeverageController) DeleteBeverage(ctx *gin.Context) {
	ID := ctx.Param("id")

	if controller.service.DeleteBeverage(ID) {
		fmt.Fprint(ctx.Writer, "{\"status\":\"OK\"}")
	} else {
		fmt.Fprint(ctx.Writer, "{\"status\":\"ERROR\"}")
	}

}
