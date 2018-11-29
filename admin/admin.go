package main

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
	"strings"
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
	Name1 string //Username optional to filter for
	Name2 string //Username optional to filter for
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
		println(strings.ToLower(cmd.Type))
		println(string(cmd.Payload))
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
		case "backup":
			bkpcmd := PerformBackupCommand{}
			json.Unmarshal(cmd.Payload, &bkpcmd)
			io.Copy(con, bytes.NewBuffer(as.doBackup(&bkpcmd)))
		default:
			println("Unknown command type")
		}
	}
}

func (as *AdminServer) doBackup(cmd *PerformBackupCommand) []byte {
	err := os.MkdirAll(cmd.Path, 0700)
	if err != nil {
		return []byte(err.Error())
	}
	//as.ur.Backup(path.Join(cmd.Path,users.bolt))
	//as.br.Backup(path.Join(cmd.Path,beverages.bolt))
	//as.ar.Backup(path.Join(cmd.Path,accounts.bolt))
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

func main() {
	perms, err := permissions.NewPermissions("../data")
	if err != nil {
		panic(err.Error())
	}
	ur, err := authStuff.NewUserRepo("../data")
	if err != nil {
		panic(err.Error())
	}
	br, err := beverages.NewBeverageRepo("../data")
	if err != nil {
		panic(err.Error())
	}
	ar, err := accounts.NewAccountRepo("../data")
	if err != nil {
		panic(err.Error())
	}
	os.Remove("./admin.socket")
	as, err := NewAdminServer("./admin.socket", ur, ar, br, perms)
	if err != nil {
		panic(err.Error())
	}
	as.StartAccepting()
}
