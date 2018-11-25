package reports

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/killingspark/hadibar/permissions"

	"github.com/gin-gonic/gin"
	"github.com/killingspark/hadibar/accounts"
	"github.com/killingspark/hadibar/authStuff"
	"github.com/killingspark/hadibar/beverages"
	"github.com/killingspark/hadibar/restapi"
	"github.com/killingspark/hadibar/settings"
)

type ReportsController struct {
	bevsrv *beverages.BeverageService
	accsrv *accounts.AccountService
}

func NewReportsController(perms *permissions.Permissions) (*ReportsController, error) {
	bevsrv, err := beverages.NewBeverageService(settings.S.DataDir, perms)
	if err != nil {
		return nil, err
	}
	accsrv, err := accounts.NewAccountService(settings.S.DataDir, perms)
	if err != nil {
		return nil, err
	}
	return &ReportsController{bevsrv: bevsrv, accsrv: accsrv}, nil
}

func makeAccountTableRow(acc *accounts.Account) string {
	return "<tr><td>" + acc.Owner.Name + "</td><td>" + strconv.Itoa(acc.Value) + "</td></tr>"
}

func (rc *ReportsController) GenerateAccountList(ctx *gin.Context) {
	var info *authStuff.LoginInfo
	if inter, ok := ctx.Get("logininfo"); ok {
		info, ok = inter.(*authStuff.LoginInfo)
		if !ok {
			response, _ := restapi.NewErrorResponse("Something went wrong while processing the username").Marshal()
			fmt.Fprint(ctx.Writer, string(response))
			ctx.Abort()
			return
		}
	}

	accs, err := rc.accsrv.GetAccounts(info.Name)
	if err != nil {
		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}
	report := "<table id='accreport'><th>Name</th><th>Value</th>"
	for _, acc := range accs {
		report += makeAccountTableRow(acc)
	}
	report += "</table>"
	response, _ := restapi.NewOkResponse(report).Marshal()
	fmt.Fprint(ctx.Writer, string(response))
	ctx.Next()
}

func makeBeverageTableRow(acc *accounts.Account, bevs int) string {
	row := "<tr><td>" + acc.Owner.Name + "</td>"
	for i := 0; i < bevs; i++ {
		row += "<td></td>"
	}
	return row + "</tr>"
}

func (rc *ReportsController) GenerateBeverageMatrix(ctx *gin.Context) {
	var info *authStuff.LoginInfo
	if inter, ok := ctx.Get("logininfo"); ok {
		info, ok = inter.(*authStuff.LoginInfo)
		if !ok {
			response, _ := restapi.NewErrorResponse("Something went wrong while processing the username").Marshal()
			fmt.Fprint(ctx.Writer, string(response))
			ctx.Abort()
			return
		}
	}

	bevs, err := rc.bevsrv.GetBeverages(info.Name)
	if err != nil {
		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}
	report := "<table id='matrix'><th></th>"
	for _, bev := range bevs {
		report += "<th>" + bev.Name + ": " + strconv.Itoa(bev.Value) + "ct </th>"
	}
	accs, err := rc.accsrv.GetAccounts(info.Name)
	if err != nil {
		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}
	for _, acc := range accs {
		report += makeBeverageTableRow(acc, len(bevs))
	}
	report += "</table>"
	response, _ := restapi.NewOkResponse(report).Marshal()
	fmt.Fprint(ctx.Writer, string(response))
	ctx.Next()
}

func makeTransactionRow(srcName, trgtName string, amount int, time string) string {
	row := "<tr><td>" + srcName + "</td><td>" + trgtName + "</td><td>" + strconv.Itoa(amount) + "</td><td>" + time + "</td></tr>"
	return row
}

func (rc *ReportsController) GenerateTransactionList(ctx *gin.Context) {
	var info *authStuff.LoginInfo
	if inter, ok := ctx.Get("logininfo"); ok {
		info, ok = inter.(*authStuff.LoginInfo)
		if !ok {
			response, _ := restapi.NewErrorResponse("Something went wrong while processing the username").Marshal()
			fmt.Fprint(ctx.Writer, string(response))
			ctx.Abort()
			return
		}
	}
	accID := ctx.PostForm("accid")

	txs, err := rc.accsrv.GetTransactions(accID, info.Name)
	if err != nil {
		response, _ := restapi.NewErrorResponse(err.Error()).Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}

	idMap := make(map[string]string)
	report := "<table><th>Source</th><th>Target</th><th>Amount</th><th>Time</th>"

	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Timestamp.After(txs[j].Timestamp)
	})

	for _, tx := range txs {
		var srcName string
		if tx.SourceID == "0" {
			srcName = "Outside"
		} else {
			var ok bool
			srcName, ok = idMap[tx.SourceID]
			if !ok {
				acc, err := rc.accsrv.GetAccount(tx.SourceID, info.Name)
				if err != nil {
					continue
				}
				idMap[tx.SourceID] = acc.Owner.Name
				srcName = acc.Owner.Name
			}
		}

		trgtName, ok := idMap[tx.TargetID]
		if !ok {
			acc, err := rc.accsrv.GetAccount(tx.TargetID, info.Name)
			if err != nil {
				println(err.Error())
				continue
			}
			idMap[tx.TargetID] = acc.Owner.Name
			trgtName = acc.Owner.Name
		}
		report += makeTransactionRow(srcName, trgtName, tx.Amount, tx.Timestamp.Format(time.UnixDate))
	}
	report += "</table>"
	response, _ := restapi.NewOkResponse(report).Marshal()
	fmt.Fprint(ctx.Writer, string(response))
	ctx.Next()
}
