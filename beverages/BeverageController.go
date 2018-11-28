package beverages

import (
	"fmt"

	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/killingspark/hadibar/authStuff"
	"github.com/killingspark/hadibar/permissions"
	"github.com/killingspark/hadibar/restapi"
)

//BeverageController : Controller for the Beverages
type BeverageController struct {
	service *BeverageService
}

//NewBeverageController creates a new BeverageController and initializes its service
func NewBeverageController(perms *permissions.Permissions, datadir string) (*BeverageController, error) {
	bc := &BeverageController{}
	var err error
	bc.service, err = NewBeverageService(datadir, perms)
	if err != nil {
		return nil, err
	}
	return bc, nil
}

//GetBeverages responds with all existing Beverages
func (controller *BeverageController) GetBeverages(ctx *gin.Context) {
	info, err := authStuff.GetLoginInfoFromCtx(ctx)
	if err != nil {
		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}

	bevs, err := controller.service.GetBeverages(info.Name)
	if err != nil {
		errResp, _ := restapi.NewErrorResponse("Couldnt get the beverage array").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	response, err := restapi.NewOkResponse(bevs).Marshal()
	if err != nil {
		errResp, _ := restapi.NewErrorResponse("Couldnt marshal the beverage array").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}
	fmt.Fprint(ctx.Writer, string(response))
	ctx.Next()
}

//GetBeverage responds with the beverage identified by the ID in the query
func (controller *BeverageController) GetBeverage(ctx *gin.Context) {
	info, err := authStuff.GetLoginInfoFromCtx(ctx)
	if err != nil {
		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}

	ID, ok := ctx.GetQuery("id")
	if !ok {
		errResp, _ := restapi.NewErrorResponse("No ID given").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	bev, err := controller.service.GetBeverage(ID, info.Name)
	if err == nil {
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
		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}
}

//NewBeverage creates a new beverage with the given form-values "value", "name" and "available" and returns it
func (controller *BeverageController) NewBeverage(ctx *gin.Context) {
	info, err := authStuff.GetLoginInfoFromCtx(ctx)
	if err != nil {
		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
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

	na, err := strconv.Atoi(ctx.PostForm("available"))
	if err != nil {
		errResp, _ := restapi.NewErrorResponse("Invalid available").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	nb, err := controller.service.NewBeverage(info.Name, ctx.PostForm("name"), nv, na)
	if err != nil {
		errResp, _ := restapi.NewErrorResponse("Couldnt save new beverage: " + err.Error()).Marshal()
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

//UpdateBeverage updates the beverage identified by the ID the query with the given form-values "value", "name" and "available" and returns it
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
	na, err := strconv.Atoi(ctx.PostForm("available"))
	if err != nil {
		errResp, _ := restapi.NewErrorResponse("Invalid available").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}
	nn := ctx.PostForm("name")

	info, err := authStuff.GetLoginInfoFromCtx(ctx)
	if err != nil {
		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}

	nb, err := controller.service.UpdateBeverage(ID, info.Name, nn, nv, na)
	if err != nil {
		errResp, _ := restapi.NewErrorResponse("Couldnt update beverage: " + err.Error()).Marshal()
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

//GivePermissionToUser gives the other user permission to read/alter the beverage
func (controller *BeverageController) GivePermissionToUser(ctx *gin.Context) {
	ID, ok := ctx.GetQuery("id")
	if !ok {
		errResp, _ := restapi.NewErrorResponse("No ID given").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	newowner := ctx.PostForm("newowner")
	if newowner == "" {
		errResp, _ := restapi.NewErrorResponse("No newowner given").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	info, err := authStuff.GetLoginInfoFromCtx(ctx)
	if err != nil {
		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}

	err = controller.service.GivePermissionToUser(ID, info.Name, newowner, permissions.CRUD)
	if err != nil {
		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}
	response, _ := restapi.NewOkResponse("").Marshal()
	fmt.Fprint(ctx.Writer, string(response))
	ctx.Next()
}

//DeleteBeverage deletes the beverage identified by the ID in the query and responds with a YEP/NOPE
func (controller *BeverageController) DeleteBeverage(ctx *gin.Context) {
	ID, ok := ctx.GetQuery("id")
	if !ok {
		errResp, _ := restapi.NewErrorResponse("No ID given").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	info, err := authStuff.GetLoginInfoFromCtx(ctx)
	if err != nil {
		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}

	if err := controller.service.DeleteBeverage(ID, info.Name); err == nil {
		response, _ := restapi.NewOkResponse("").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Next()
	} else {
		response, _ := restapi.NewErrorResponse("Coulnt delete the beverage: " + err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}
}
