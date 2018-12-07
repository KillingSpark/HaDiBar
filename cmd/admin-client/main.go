package main

import (
	"bytes"
	"encoding/json"
	"github.com/killingspark/hadibar/admin"
	"net"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app    = kingpin.New("hadibar-admin", "A client for the admin-interface of hadibar")
	socket = app.Flag("socket", "The location of the socket file").Default("./control.socket").Short('s').String()

	delusr     = app.Command("delusr", "Delete a user")
	delusrName = delusr.Arg("name", "Name of user").Required().String()

	clean  = app.Command("clean", "Clean orphaned accounts and beverages")
	lsusrs = app.Command("lsusrs", "List all user")

	lsaccs        = app.Command("lsaccs", "List all accounts (and filter for a specific user)")
	lsaccsusrName = lsaccs.Arg("name", "Name of the user for filtering").String()

	lsbevs        = app.Command("lsbevs", "List all beverages (and filter for a specific user)")
	lsbevsusrName = lsbevs.Arg("name", "Name of the user for filtering").String()

	lstxs  = app.Command("lstxs", "List all transactions (and filter for a specific accounts)")
	txsID1 = lstxs.Arg("id1", "ID of one account for filtering").String()
	txsID2 = lstxs.Arg("id2", "ID of other account for filtering").String()
)

type cmdWrap struct {
	Type string      `json:"type"`
	Cmd  interface{} `json:"cmd"`
}

func main() {
	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	var err error
	con, err = net.Dial("unix", *socket)
	if err != nil {
		panic(err.Error())
	}

	var cmd interface{}

	switch command {
	// delete user
	case delusr.FullCommand():
		cmd = cmdWrap{Type: "deleteuser", Cmd: admin.DeleteUserCommand{Name: *delusrName}}
	case clean.FullCommand():
		cmd = cmdWrap{Type: "clean"}
	case lsusrs.FullCommand():
		cmd = cmdWrap{Type: "listusers", Cmd: admin.ListUsersCommand{}}
	case lsaccs.FullCommand():
		cmd = cmdWrap{Type: "listaccs", Cmd: admin.ListAccountsCommand{Name: *lsaccsusrName}}
	case lsbevs.FullCommand():
		cmd = cmdWrap{Type: "listbevs", Cmd: admin.ListBeveragesCommand{Name: *lsbevsusrName}}
	case lstxs.FullCommand():
		cmd = cmdWrap{Type: "listtxs", Cmd: admin.ListTransactionssCommand{ID1: *txsID1, ID2: *txsID2}}
	default:
		println("unknown: " + command)
		return
	}
	b, err := json.Marshal(cmd)
	if err != nil {
		panic(err.Error())
	}
	err = writeCommand(b)
	if err != nil {
		panic(err.Error())
	}
}

var con net.Conn
var buf = make([]byte, 4069)
var pretty = bytes.NewBuffer(make([]byte, 4069))

func writeCommand(line []byte) error {
	xs := 0
	for xs < len(line) {
		x, err := con.Write(line)
		xs += x
		if err != nil {
			return err
		}
	}

	c, err := con.Read(buf)
	if err != nil {
		return err
	}
	json.Indent(pretty, buf[:c], "", "  ")
	println(string(pretty.Bytes()))
	pretty.Reset()
	return nil
}
