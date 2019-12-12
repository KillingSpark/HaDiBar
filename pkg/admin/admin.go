package admin

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/killingspark/hadibar/pkg/accounts"
	"github.com/killingspark/hadibar/pkg/authStuff"
	"github.com/killingspark/hadibar/pkg/beverages"
	"github.com/killingspark/hadibar/pkg/permissions"
)

type AdminServer struct {
	lstn       net.Listener
	ur         *authStuff.UserRepo
	br         *beverages.BeverageRepo
	ar         *accounts.AccountRepo
	perm       *permissions.Permissions
	socketPath string
}

// `{"Result": "Err", "Text":"` + err.Error() + `"}`
type ErrorResponse struct {
	Result string `json:"Result"`
	Text   string `json:"Text"`
}

func NewUnixAdminServer(pathToSocket string, usrRepo *authStuff.UserRepo, accRepo *accounts.AccountRepo, bevRepo *beverages.BeverageRepo, perms *permissions.Permissions) (*AdminServer, error) {
	as := &AdminServer{}
	var err error
	as.lstn, err = net.Listen("unix", pathToSocket)
	if err != nil {
		return nil, err
	}
	as.ur = usrRepo
	as.ar = accRepo
	as.br = bevRepo
	as.perm = perms
	as.socketPath = pathToSocket
	return as, nil
}

func NewTlsAdminServer(addr, certPath, keyPath, caCertPath string, tlsClientcertRequired bool, usrRepo *authStuff.UserRepo, accRepo *accounts.AccountRepo, bevRepo *beverages.BeverageRepo, perms *permissions.Permissions) (*AdminServer, error) {
	as := &AdminServer{}
	var err error

	tcpLstn, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	conf := &tls.Config{}
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	conf.Certificates = append(conf.Certificates, cert)

	if caCertPath != "" {
		certPEMBlock, err := ioutil.ReadFile(caCertPath)
		if err != nil {
			return nil, err
		}
		conf.RootCAs = x509.NewCertPool()
		if !conf.RootCAs.AppendCertsFromPEM(certPEMBlock) {
			panic("Couldnt load CA pem")
		}
		conf.ClientCAs = conf.RootCAs
	}

	if tlsClientcertRequired {
		conf.ClientAuth = tls.RequireAndVerifyClientCert
	}

	as.lstn = tls.NewListener(tcpLstn, conf)

	as.ur = usrRepo
	as.ar = accRepo
	as.br = bevRepo
	as.perm = perms
	return as, nil
}

func NewTcpAdminServer(addr string, usrRepo *authStuff.UserRepo, accRepo *accounts.AccountRepo, bevRepo *beverages.BeverageRepo, perms *permissions.Permissions) (*AdminServer, error) {
	as := &AdminServer{}
	var err error

	as.lstn, err = net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	as.ur = usrRepo
	as.ar = accRepo
	as.br = bevRepo
	as.perm = perms
	return as, nil
}

func (as *AdminServer) Close() error {
	if as.socketPath != "" {
		return os.Remove(as.socketPath)
	} else {
		return as.lstn.Close()
	}
}

type Command struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"cmd"`
}

type ListUsersCommand struct {
}

type DeleteUserCommand struct {
	Name string
}

type ListBeveragesCommand struct {
	Name string //Username optional to filter for
}

type ListAccountsCommand struct {
	Name string //Username optional to filter for
}

type ListTransactionssCommand struct {
	ID1  string    //Username optional to filter for
	ID2  string    //Username optional to filter for
	From time.Time //Oldest transaction
	To   time.Time //Newest transaction
}

type PerformBackupCommand struct {
	Path string
}

func (as *AdminServer) StartAccepting() {
	for {
		con, err := as.lstn.Accept()
		if err != nil {
			continue
		}
		go as.handleCon(con)
	}
}

func (as *AdminServer) handleCon(con net.Conn) {
	dec := json.NewDecoder(con)
	cmd := Command{}
	err := dec.Decode(&cmd)
	if err != nil {
		con.Close()
		return
	}
	switch strings.ToLower(cmd.Type) {
	case "listusers":
		lucmd := ListUsersCommand{}
		io.Copy(con, bytes.NewBuffer(as.listUsers(&lucmd)))
	case "deleteuser":
		ducmd := DeleteUserCommand{}
		json.Unmarshal(cmd.Payload, &ducmd)
		io.Copy(con, bytes.NewBuffer(as.deleteUser(&ducmd)))
	case "listbevs":
		lbcmd := ListBeveragesCommand{}
		json.Unmarshal(cmd.Payload, &lbcmd)
		io.Copy(con, bytes.NewBuffer(as.listBevs(&lbcmd)))
	case "listaccs":
		lacmd := ListAccountsCommand{}
		json.Unmarshal(cmd.Payload, &lacmd)
		io.Copy(con, bytes.NewBuffer(as.listAccs(&lacmd)))
	case "listtxs":
		ltxscmd := ListTransactionssCommand{}
		json.Unmarshal(cmd.Payload, &ltxscmd)
		io.Copy(con, bytes.NewBuffer(as.listTransactions(&ltxscmd)))
	case "backup":
		bkpcmd := PerformBackupCommand{}
		json.Unmarshal(cmd.Payload, &bkpcmd)
		io.Copy(con, bytes.NewBuffer(as.doBackup(&bkpcmd)))
	case "clean":
		io.Copy(con, bytes.NewBuffer(as.cleanUpOrphaned()))
	default:
		log.WithFields(log.Fields{"cmd": strings.ToLower(cmd.Type)}).Warn("Unknown command")
	}
	con.Close()

}

func (as *AdminServer) listTransactions(cmd *ListTransactionssCommand) []byte {
	log.WithFields(log.Fields{"cmd": "listtxs", "from": cmd.From.String(), "to": cmd.To.String(), "ID1": cmd.ID1, "ID2": cmd.ID2}).Debug("Command received")

	txs, err := as.ar.GetTransactions()
	if err != nil {
		result := ErrorResponse{Result: "Err", Text: err.Error()}
		marshalled, _ := json.Marshal(&result)
		return marshalled
	}
	txsFiltered := make([]*accounts.Transaction, 0)

	for _, tx := range txs {
		if tx.SourceID == cmd.ID1 || tx.TargetID == cmd.ID1 ||
			tx.SourceID == cmd.ID2 || tx.TargetID == cmd.ID2 || (cmd.ID1 == "" && cmd.ID2 == "") {
			if cmd.To.Equal(cmd.From) || (tx.Timestamp.Before(cmd.To) && tx.Timestamp.After(cmd.From)) {
				txsFiltered = append(txsFiltered, tx)
			}
		}
	}
	marshed, err := json.Marshal(txsFiltered)
	if err != nil {
		return []byte(`{"Result": "Err", "Text":"` + err.Error() + `"}`)
	}
	return marshed
}

func (as *AdminServer) doBackup(cmd *PerformBackupCommand) []byte {
	log.WithFields(log.Fields{"cmd": "backup", "path": cmd.Path}).Debug("Command received")

	bkpPath := cmd.Path
	if bkpPath == "" {
		result := ErrorResponse{Result: "Err", Text: "No path given"}
		marshalled, _ := json.Marshal(&result)
		return marshalled
	}
	err := os.MkdirAll(bkpPath, 0700)
	if err != nil {
		result := ErrorResponse{Result: "Err", Text: err.Error()}
		marshalled, _ := json.Marshal(&result)
		return marshalled
	}

	//Lock all repos before doing anything
	as.ur.Lock.RLock()
	defer as.ur.Lock.RUnlock()
	as.br.Lock.RLock()
	defer as.ur.Lock.RUnlock()
	as.ar.Lock.RLock()
	defer as.ur.Lock.RUnlock()

	err = as.ur.BackupTo(path.Join(bkpPath, "users.bolt"))
	if err != nil {
		return []byte(`{"Result": "Err", "Text":"` + err.Error() + `"}`)
	}
	err = as.br.BackupTo(path.Join(bkpPath, "beverages.bolt"))
	if err != nil {
		return []byte(`{"Result": "Err", "Text":"` + err.Error() + `"}`)
	}
	err = as.ar.BackupTo(path.Join(bkpPath, "accounts.bolt"))
	if err != nil {
		return []byte(`{"Result": "Err", "Text":"` + err.Error() + `"}`)
	}
	err = as.perm.BackupTo(path.Join(bkpPath, "permissions.bolt"))
	if err != nil {
		return []byte(`{"Result": "Err", "Text":"` + err.Error() + `"}`)
	}
	return []byte(`{"Result": "OK"}`)
}

func (as *AdminServer) cleanUpOrphaned() []byte {
	log.WithFields(log.Fields{"cmd": "clean"}).Debug("Command received")

	perms, err := as.perm.GetAllAsMap()
	if err != nil {
		result := ErrorResponse{Result: "Err", Text: err.Error()}
		marshalled, _ := json.Marshal(&result)
		return marshalled
	}
	objsFound := make(map[string]bool)
	for _, usrobjs := range perms {
		for objID := range usrobjs {
			objsFound[objID] = true
		}
	}

	accs, err := as.ar.GetAllAccounts()
	if err != nil {
		result := ErrorResponse{Result: "Err", Text: err.Error()}
		marshalled, _ := json.Marshal(&result)
		return marshalled
	}
	bevs, err := as.br.GetAllBeverages()
	if err != nil {
		result := ErrorResponse{Result: "Err", Text: err.Error()}
		marshalled, _ := json.Marshal(&result)
		return marshalled
	}
	deleted := 0
	for _, acc := range accs {
		if _, ok := objsFound[acc.ID]; !ok {
			as.ar.DeleteInstance(acc.ID)
			deleted++
		}
	}
	for _, bev := range bevs {
		if _, ok := objsFound[bev.ID]; !ok {
			as.br.DeleteInstance(bev.ID)
			deleted++
		}
	}
	return []byte(`{"Result": "OK", "Text":  "cleaned ` + strconv.Itoa(deleted) + `"}`)
}

func (as *AdminServer) deleteUser(cmd *DeleteUserCommand) []byte {
	log.WithFields(log.Fields{"cmd": "deleteuser", "name": cmd.Name}).Debug("Command received")

	err := as.ur.DeleteInstance(cmd.Name)
	if err != nil {
		result := ErrorResponse{Result: "Err", Text: err.Error()}
		marshalled, _ := json.Marshal(&result)
		return marshalled
	}
	as.perm.RemoveUsersPermissions(cmd.Name)
	return []byte(`{"Result": "OK"}`)
}

func (as *AdminServer) listUsers(cmd *ListUsersCommand) []byte {
	log.WithFields(log.Fields{"cmd": "listusers"}).Debug("Command received")

	users, err := as.ur.GetAllUsers()
	if err != nil {
		result := ErrorResponse{Result: "Err", Text: err.Error()}
		marshalled, _ := json.Marshal(&result)
		return marshalled
	}
	marshed, err := json.Marshal(users)
	if err != nil {
		result := ErrorResponse{Result: "Err", Text: err.Error()}
		marshalled, _ := json.Marshal(&result)
		return marshalled
	}
	return marshed
}

func (as *AdminServer) listAccs(cmd *ListAccountsCommand) []byte {
	log.WithFields(log.Fields{"cmd": "listaccs"}).Debug("Command received")

	accs, err := as.ar.GetAllAccounts()
	if err != nil {
		result := ErrorResponse{Result: "Err", Text: err.Error()}
		marshalled, _ := json.Marshal(&result)
		return marshalled
	}

	var accsFiltered []*accounts.Account

	//filter for name if given
	if cmd.Name != "" {
		for _, acc := range accs {
			ok, err := as.perm.CheckPermissionAny(acc.ID, cmd.Name, permissions.Read, permissions.Update, permissions.Delete, permissions.CRUD)
			if err != nil {
				continue
			}
			if ok {
				accsFiltered = append(accsFiltered, acc)
			}
		}
	} else {
		accsFiltered = accs
	}
	marshed, err := json.Marshal(accsFiltered)
	if err != nil {
		result := ErrorResponse{Result: "Err", Text: err.Error()}
		marshalled, _ := json.Marshal(&result)
		return marshalled
	}
	return marshed
}

func (as *AdminServer) listBevs(cmd *ListBeveragesCommand) []byte {
	log.WithFields(log.Fields{"cmd": "listbevs"}).Debug("Command received")

	bevs, err := as.br.GetAllBeverages()
	if err != nil {
		result := ErrorResponse{Result: "Err", Text: err.Error()}
		marshalled, _ := json.Marshal(&result)
		return marshalled
	}

	var bevsFiltered []*beverages.Beverage

	//filter for name if given
	if cmd.Name != "" {
		for _, bev := range bevs {
			ok, err := as.perm.CheckPermissionAny(bev.ID, cmd.Name, permissions.Read, permissions.Update, permissions.Delete, permissions.CRUD)
			if err != nil {
				continue
			}
			if ok {
				bevsFiltered = append(bevsFiltered, bev)
			}
		}
	} else {
		bevsFiltered = bevs
	}
	marshed, err := json.Marshal(bevsFiltered)
	if err != nil {
		result := ErrorResponse{Result: "Err", Text: err.Error()}
		marshalled, _ := json.Marshal(&result)
		return marshalled
	}
	return marshed
}
