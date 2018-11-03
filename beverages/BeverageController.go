package beverages

import (
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
	response, err := restapi.NewOkResponse(controller.service.GetBeverages()).Marshal()
	if err != nil {
		errResp, _ := restapi.NewErrorResponse("Couldnt marshal the beverage array").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}
	fmt.Fprint(ctx.Writer, string(response))
	ctx.Next()
}

//GetBeverage responds with the beverage identified by beverage/:id
func (controller *BeverageController) GetBeverage(ctx *gin.Context) {
	ID, ok := ctx.GetQuery("id")
	if !ok {
		errResp, _ := restapi.NewErrorResponse("No ID given").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	bev, ok := controller.service.GetBeverage(ID)
	if ok {
		response, err := restapi.NewOkResponse(bev).Marshal()
		if err != nil {
			errResp, _ := restapi.NewErrorResponse("Couldnt marshal the beverage object").Marshal()
			fmt.Fprint(ctx.Writer, string(errResp))
			ctx.Abort()
			return
		}
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Next()
	} else {
		response, _ := restapi.NewErrorResponse("").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}
}

//NewBeverage creates a new beverage with the given form-values "value" and "name" and returns it
func (controller *BeverageController) NewBeverage(ctx *gin.Context) {
	nv, err := strconv.Atoi(ctx.PostForm("value"))
	if err != nil {
		errResp, _ := restapi.NewErrorResponse("Invalid value").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	nb, ok := controller.service.NewBeverage(ctx.PostForm("name"), nv)
	if !ok {
		errResp, _ := restapi.NewErrorResponse("Couldnt save new beverage").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}
	response, err := restapi.NewOkResponse(nb).Marshal()
	if err != nil {
		errResp, _ := restapi.NewErrorResponse("Couldnt marshal the beverage object").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}
	fmt.Fprint(ctx.Writer, string(response))
	ctx.Next()
}

//UpdateBeverage updates the beverage identified by /beverage/:id with the given form-values "value" and "name" and returns it
func (controller *BeverageController) UpdateBeverage(ctx *gin.Context) {
	ID, ok := ctx.GetQuery("id")
	if !ok {
		errResp, _ := restapi.NewErrorResponse("No ID given").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	nv, err := strconv.Atoi(ctx.PostForm("value"))
	if err != nil {
		errResp, _ := restapi.NewErrorResponse("Invalid value").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}
	nn := ctx.PostForm("name")
	nb, ok := controller.service.UpdateBeverage(ID, nn, nv)
	if !ok {
		errResp, _ := restapi.NewErrorResponse("Couldnt update beverage").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	response, err := restapi.NewOkResponse(nb).Marshal()
	if err != nil {
		errResp, _ := restapi.NewErrorResponse("Couldnt marshal the beverage object").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}
	fmt.Fprint(ctx.Writer, string(response))
	ctx.Next()
}

//DeleteBeverage deletes the beverage identified by /beverage/:id and responds with a YEP/NOPE
func (controller *BeverageController) DeleteBeverage(ctx *gin.Context) {
	ID, ok := ctx.GetQuery("id")
	if !ok {
		errResp, _ := restapi.NewErrorResponse("No ID given").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	if controller.service.DeleteBeverage(ID) {
		response, _ := restapi.NewOkResponse("").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Next()
	} else {
		response, _ := restapi.NewErrorResponse("").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}
}
