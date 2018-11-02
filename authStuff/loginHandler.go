package authStuff

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"hash"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/sha3"
)

type jsonUserDatabase struct {
	path   string //where is the file with the users
	users  map[string]*LoginInfo
	hasher hash.Hash
}

func NewJsonUserDatabase() *jsonUserDatabase {
	db := &jsonUserDatabase{}
	db.hasher = sha3.New256()
	db.path = os.ExpandEnv("$HOME") + "/.cache/hadibarusers"

	return db
}

func (db *jsonUserDatabase) load() error {
	jsonFile, err := os.Open(db.path)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(byteValue), &db.users)
	if err != nil {
		return err
	}

	return nil
}

func (db *jsonUserDatabase) save() error {
	jsonFile, err := os.OpenFile(db.path, os.O_RDWR, 0)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	enc, err := json.Marshal(db.users)
	if err != nil {
		return err
	}

	_, err = jsonFile.Write(enc)
	if err != nil {
		return err
	}

	return nil
}

var ErrUserNotKnown = errors.New("User not in database")
var ErrWrongCredetials = errors.New("Wrong creds for username")

func saltPw(hasher hash.Hash, pw, salt string) string {
	hasher.Reset()
	hasher.Write([]byte(pw + salt))
	saltedpw := make([]byte, 1024)

	base64.StdEncoding.Encode(saltedpw, hasher.Sum(nil))
	return string(saltedpw)
}

func (db *jsonUserDatabase) isValid(user, pw string) (*LoginInfo, error) {
	err := db.load()
	if err != nil {
		return nil, err
	}

	lgi, ok := db.users[user]
	if !ok {
		return nil, ErrUserNotKnown
	}
	if saltPw(db.hasher, pw, lgi.Salt) == lgi.Pwhash {
		return lgi, nil
	}
	return nil, ErrWrongCredetials
}
