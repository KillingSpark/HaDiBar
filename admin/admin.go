package admin

import (
	"bytes"
	"encoding/json"
	"github.com/killingspark/hadibar/accounts"
	"github.com/killingspark/hadibar/authStuff"
	"github.com/killingspark/hadibar/beverages"
	"github.com/killingspark/hadibar/permissions"
	"io"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

type AdminServer struct {
	lstn       net.Listener
	ur         *authStuff.UserRepo
	br         *beverages.BeverageRepo
	ar         *accounts.AccountRepo
	perm       *permissions.Permissions
	socketPath string
}

func NewAdminServer(pathToSocket string, usrRepo *authStuff.UserRepo, accRepo *accounts.AccountRepo, bevRepo *beverages.BeverageRepo, perms *permissions.Permissions) (*AdminServer, error) {
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

func (as *AdminServer) Close() error {
	return os.Remove(as.socketPath)
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
	for {
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
		default:
			println("Unknown command type received on the unix socket: " + cmd.Type)
		}
	}
}

func (as *AdminServer) listTransactions(cmd *ListTransactionssCommand) []byte {
	txs, err := as.ar.GetTransactions()
	if err != nil {
		return []byte(err.Error())
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
		return []byte(err.Error())
	}
	return marshed
}

func (as *AdminServer) doBackup(cmd *PerformBackupCommand) []byte {
	bkpPath := cmd.Path
	if bkpPath == "" {
		return []byte("No path given")
	}
	err := os.MkdirAll(bkpPath, 0700)
	if err != nil {
		return []byte(err.Error())
	}
	err = as.ur.BackupTo(path.Join(bkpPath, "users.bolt"))
	if err != nil {
		return []byte(err.Error())
	}
	err = as.br.BackupTo(path.Join(bkpPath, "beverages.bolt"))
	if err != nil {
		return []byte(err.Error())
	}
	err = as.ar.BackupTo(path.Join(bkpPath, "accounts.bolt"))
	if err != nil {
		return []byte(err.Error())
	}
	err = as.perm.BackupTo(path.Join(bkpPath, "permissions.bolt"))
	if err != nil {
		return []byte(err.Error())
	}
	return []byte("OK")
}

func (as *AdminServer) deleteUser(cmd *DeleteUserCommand) []byte {
	err := as.ur.DeleteInstance(cmd.Name)
	if err != nil {
		return []byte(err.Error())
	}
	return []byte("OK")
}

func (as *AdminServer) listUsers(cmd *ListUsersCommand) []byte {
	users, err := as.ur.GetAllUsers()
	if err != nil {
		return []byte(err.Error())
	}
	marshed, err := json.Marshal(users)
	if err != nil {
		return []byte(err.Error())
	}
	return marshed
}

func (as *AdminServer) listAccs(cmd *ListAccountsCommand) []byte {
	accs, err := as.ar.GetAllAccounts()
	if err != nil {
		return []byte(err.Error())
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
		return []byte(err.Error())
	}
	return marshed
}

func (as *AdminServer) listBevs(cmd *ListBeveragesCommand) []byte {
	bevs, err := as.br.GetAllBeverages()
	if err != nil {
		return []byte(err.Error())
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
		return []byte(err.Error())
	}
	return marshed
}
