package authStuff

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"hash"
	"io/ioutil"
	"math"
	"os"

	"github.com/killingspark/HaDiBar/settings"

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
	db.path = os.ExpandEnv(settings.S.UserPath)

	return db
}

func (db *jsonUserDatabase) Add(new *LoginInfo) error {
	if _, ok := db.users[new.Name]; ok {
		return errors.New("already exists")
	}
	db.users[new.Name] = new
	return nil
}

func (db *jsonUserDatabase) Load() error {
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

func (db *jsonUserDatabase) Save() error {
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

func SaltPw(hasher hash.Hash, pw, salt string) string {
	hasher.Reset()
	hasher.Write([]byte(pw + salt))
	saltedpw := make([]byte, 4*int(math.Ceil((float64(32)/3))))

	base64.StdEncoding.Encode(saltedpw, hasher.Sum(nil))
	return string(saltedpw)
}

func (db *jsonUserDatabase) isValid(user, pw string) (*LoginInfo, error) {
	err := db.Load()
	if err != nil {
		return nil, err
	}

	lgi, ok := db.users[user]
	if !ok {
		return nil, ErrUserNotKnown
	}
	if SaltPw(db.hasher, pw, lgi.Salt) == lgi.Pwhash {
		return lgi, nil
	}
	return nil, ErrWrongCredetials
}
