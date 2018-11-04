package authStuff

import (
	"encoding/base64"
	"errors"
	"hash"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/nanobox-io/golang-scribble"

	"golang.org/x/crypto/sha3"
)

var (
	collectionName = "user"
)

type LoginService struct {
	userRepo *scribble.Driver
	hasher   hash.Hash
}

func NewLoginService(path string) (*LoginService, error) {
	ls := &LoginService{}
	ls.hasher = sha3.New256()
	var err error
	ls.userRepo, err = scribble.New(path, nil)

	if err != nil {
		return nil, err
	}

	return ls, nil
}

var ErrUsernameTaken = errors.New("already exists")
var ErrUserNotKnown = errors.New("User not in database")
var ErrWrongCredetials = errors.New("Wrong creds for username")

func (ls *LoginService) Add(new *LoginInfo) error {
	var user *LoginInfo
	if ls.userRepo.Read(collectionName, new.Name, user); user != nil && user.Name == new.Name {
		return ErrUsernameTaken
	}
	if err := ls.userRepo.Write(collectionName, new.Name, new); err != nil {
		return err
	}
	return nil
}

func createNewUser(hasher hash.Hash, username, passwd string) *LoginInfo {
	user := &LoginInfo{}
	user.Name = username
	user.GroupID = strconv.FormatInt(time.Now().UnixNano(), 10)
	user.Salt = saltPw(hasher, strconv.FormatInt(time.Now().UnixNano()%rand.Int63(), 10), username)
	user.Pwhash = saltPw(hasher, passwd, user.Salt)
	return user
}

func (ls *LoginService) isValid(userName, passwd string) (*LoginInfo, error) {
	var user *LoginInfo
	if err := ls.userRepo.Read(collectionName, userName, user); err != nil {
		//add unknown user with a unique groupid
		user = createNewUser(ls.hasher, userName, passwd)
		err = ls.Add(user)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	if saltPw(ls.hasher, passwd, user.Salt) == user.Pwhash {
		return user, nil
	}
	return nil, ErrWrongCredetials
}

func saltPw(hasher hash.Hash, pw, salt string) string {
	hasher.Reset()
	hasher.Write([]byte(pw + salt))
	saltedpw := make([]byte, 4*int(math.Ceil((float64(32)/3))))

	base64.StdEncoding.Encode(saltedpw, hasher.Sum(nil))
	return string(saltedpw)
}
