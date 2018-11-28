package authStuff

import (
	"encoding/json"
	scribble "github.com/nanobox-io/golang-scribble"
)

type UserRepo struct {
	db *scribble.Driver
}

var collectionName = "users"

func NewUserRepo(path string) (*UserRepo, error) {
	ar := &UserRepo{}
	var err error
	ar.db, err = scribble.New(path, nil)
	if err != nil {
		return nil, err
	}
	return ar, nil
}

func (ar *UserRepo) GetAllUsers() ([]*LoginInfo, error) {
	list, err := ar.db.ReadAll(collectionName)
	if err != nil {
		return nil, err
	}

	var res []*LoginInfo
	for _, item := range list {
		usr := &LoginInfo{}
		err := json.Unmarshal([]byte(item), usr)
		if err != nil {
			continue //skip invalied entries. maybe implement cleanup...
		}
		res = append(res, usr)
	}
	return res, nil
}

func (ar *UserRepo) SaveInstance(usr *LoginInfo) error {
	if err := ar.db.Write(collectionName, usr.Name, usr); err != nil {
		return err
	}
	return nil
}

func (ar *UserRepo) GetInstance(usrID string) (*LoginInfo, error) {
	var usr LoginInfo
	if err := ar.db.Read(collectionName, usrID, &usr); err != nil {
		return nil, err
	}
	return &usr, nil
}

func (ar *UserRepo) DeleteInstance(usrID string) error {
	return ar.db.Delete(collectionName, usrID)
}
