package beverages

import (
	"fmt"

	"strconv"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"

	"github.com/killingspark/hadibar/src/authStuff"
	"github.com/killingspark/hadibar/src/permissions"
	"github.com/killingspark/hadibar/src/restapi"
)

//BeverageController : Controller for the Beverages
type BeverageController struct {
	service *BeverageService
}

//NewBeverageController creates a new BeverageController and initializes its service
func NewBeverageController(bevService *BeverageService) *BeverageController {
	bc := &BeverageController{}
	bc.service = bevService
	return bc
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
		log.WithFields(log.Fields{"user": info.Name}).WithError(err).Error("Beverage Error GetAll")

		errResp, _ := restapi.NewErrorResponse("Couldnt get the beverage array").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	response, _ := restapi.NewOkResponse(bevs).Marshal()
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
		log.WithFields(log.Fields{"URL": ctx.Request.URL.String()}).Warn("No ID found in query")

		errResp, _ := restapi.NewErrorResponse("No ID given").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	bev, err := controller.service.GetBeverage(ID, info.Name)
	if err == nil {
		response, _ := restapi.NewOkResponse(bev).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Next()
	} else {
		log.WithFields(log.Fields{"user": info.Name}).WithError(err).Error("Beverage Error Get")

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
		log.WithFields(log.Fields{"URL": ctx.Request.URL.String()}).Warn("No int value found in postform")

		errResp, _ := restapi.NewErrorResponse("Invalid value").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	na, err := strconv.Atoi(ctx.PostForm("available"))
	if err != nil {
		log.WithFields(log.Fields{"URL": ctx.Request.URL.String()}).Warn("No int \"available\" found in postform")

		errResp, _ := restapi.NewErrorResponse("Invalid available").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	nb, err := controller.service.NewBeverage(info.Name, ctx.PostForm("name"), nv, na)
	if err != nil {
		log.WithFields(log.Fields{"user": info.Name}).WithError(err).Error("Beverage Error New")

		errResp, _ := restapi.NewErrorResponse("Couldnt save new beverage: " + err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}
	response, _ := restapi.NewOkResponse(nb).Marshal()
	fmt.Fprint(ctx.Writer, string(response))
	ctx.Next()
}

//UpdateBeverage updates the beverage identified by the ID the query with the given form-values "value", "name" and "available" and returns it
func (controller *BeverageController) UpdateBeverage(ctx *gin.Context) {
	ID, ok := ctx.GetQuery("id")
	if !ok {
		log.WithFields(log.Fields{"URL": ctx.Request.URL.String()}).Warn("No ID found in query")

		errResp, _ := restapi.NewErrorResponse("No ID given").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	nv, err := strconv.Atoi(ctx.PostForm("value"))
	if err != nil {
		log.WithFields(log.Fields{"URL": ctx.Request.URL.String()}).Warn("No int value found in postform")

		errResp, _ := restapi.NewErrorResponse("Invalid value").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}
	na, err := strconv.Atoi(ctx.PostForm("available"))
	if err != nil {
		log.WithFields(log.Fields{"URL": ctx.Request.URL.String()}).Warn("No int \"available\" found in postform")

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
		log.WithFields(log.Fields{"user": info.Name}).WithError(err).Error("Beverage Error Update")

		errResp, _ := restapi.NewErrorResponse("Couldnt update beverage: " + err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	response, _ := restapi.NewOkResponse(nb).Marshal()
	fmt.Fprint(ctx.Writer, string(response))
	ctx.Next()
}

//GivePermissionToUser gives the other user permission to read/alter the beverage
func (controller *BeverageController) GivePermissionToUser(ctx *gin.Context) {
	ID, ok := ctx.GetQuery("id")
	if !ok {
		log.WithFields(log.Fields{"URL": ctx.Request.URL.String()}).Warn("No ID found in query")

		errResp, _ := restapi.NewErrorResponse("No ID given").Marshal()
		fmt.Fprint(ctx.Writer, string(errResp))
		ctx.Abort()
		return
	}

	newowner := ctx.PostForm("newowner")
	if newowner == "" {
		log.WithFields(log.Fields{"URL": ctx.Request.URL.String()}).Warn("No newowner found in postform")

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
		log.WithFields(log.Fields{"user": info.Name}).WithError(err).Error("Beverage Error GivePermission")

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
		log.WithFields(log.Fields{"URL": ctx.Request.URL.String()}).Warn("No ID found in query")

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
		log.WithFields(log.Fields{"user": info.Name}).WithError(err).Error("Beverage Error Delete")

		response, _ := restapi.NewErrorResponse("Coulnt delete the beverage: " + err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}
}
