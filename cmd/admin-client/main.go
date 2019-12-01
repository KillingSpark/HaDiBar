package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"net"
	"os"

	"github.com/killingspark/hadibar/pkg/admin"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app        = kingpin.New("hadibar-admin", "A client for the admin-interface of hadibar")
	socket     = app.Flag("socket", "The location of the socket file").Short('s').String()
	tcpaddr    = app.Flag("tcpaddr", "The tcp address of the admin-server").Default(":8081").Short('t').String()
	useTLS     = app.Flag("tls", "use tls").Default("false").Bool()
	cert       = app.Flag("cert", "Certificate to use to identify client against the server").String()
	key        = app.Flag("key", "Key to use to identify client against the server").String()
	cacert     = app.Flag("cacert", "Public cert for ca to check the servers certificate").String()
	servername = app.Flag("servername", "Servername if none given in tcpaddr").String()

	delusr     = app.Command("delusr", "Delete a user")
	delusrName = delusr.Arg("name", "Name of user").Required().String()

	clean  = app.Command("clean", "Clean orphaned accounts and beverages")
	lsusrs = app.Command("lsusrs", "List all user")

	bckup     = app.Command("backup", "Backup databases")
	bckuppath = bckup.Arg("dir", "Directory where the databases will be saved").String()

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

	if *socket != "" {
		con, err = net.Dial("unix", *socket)
	} else {
		if *useTLS {
			conf := &tls.Config{}
			if *cacert == "" {
				conf.InsecureSkipVerify = true
			} else {
				certPEMBlock, err := ioutil.ReadFile(*cacert)
				if err != nil {
					panic(err.Error())
				}
				conf.RootCAs = x509.NewCertPool()
				if !conf.RootCAs.AppendCertsFromPEM(certPEMBlock) {
					panic("Couldnt load CA pem")
				}
				conf.ServerName = *servername
			}
			cert, err := tls.LoadX509KeyPair(*cert, *key)
			if err != nil {
				panic(err.Error())
			}
			conf.Certificates = []tls.Certificate{cert}

			con, err = tls.Dial("tcp", *tcpaddr, conf)
			if err != nil {
				panic(err.Error())
			}
		} else {
			con, err = net.Dial("tcp", *tcpaddr)
		}
	}
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
	case bckup.FullCommand():
		cmd = cmdWrap{Type: "backup", Cmd: admin.PerformBackupCommand{Path: *bckuppath}}
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

	content, err := ioutil.ReadAll(con)
	if err != nil {
		panic(err.Error())
	}

	if len(content) > 2 {
		response := admin.ErrorResponse{}
		err = json.Unmarshal(content, &response)
		if err != nil {
			panic(err.Error())
		}
		if response.Result == "Err" {
			panic("Error from server: " + response.Text)
		}
	}
	json.Indent(pretty, content, "", "  ")
	println(string(pretty.Bytes()))
	pretty.Reset()
	return nil
}
