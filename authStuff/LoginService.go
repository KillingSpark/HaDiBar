package authStuff

import (
	"encoding/base64"
	"errors"
	"hash"
	"math"
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/crypto/sha3"
)

type LoginService struct {
	userRepo *UserRepo
	hasher   hash.Hash
}

func NewLoginService(path string) (*LoginService, error) {
	ls := &LoginService{}
	ls.hasher = sha3.New256()
	var err error
	ls.userRepo, err = NewUserRepo(path)

	if err != nil {
		return nil, err
	}

	return ls, nil
}

var ErrUsernameTaken = errors.New("already exists")
var ErrUserNotKnown = errors.New("User not in database")
var ErrWrongCredetials = errors.New("Wrong creds for username")

func (ls *LoginService) Add(new *LoginInfo) error {
	user, err := ls.userRepo.GetInstance(new.Name)
	if err != nil {
		return err
	}
	if user != nil && user.Name == new.Name {
		return ErrUsernameTaken
	}
	if err := ls.userRepo.SaveInstance(new); err != nil {
		return err
	}
	return nil
}

func createNewUser(hasher hash.Hash, username, passwd string) *LoginInfo {
	user := &LoginInfo{}
	user.Name = username
	user.Salt = saltPw(hasher, strconv.FormatInt(time.Now().UnixNano()%rand.Int63(), 10), username)
	user.Pwhash = saltPw(hasher, passwd, user.Salt)
	return user
}

func (ls *LoginService) isValid(userName, passwd string) (*LoginInfo, error) {
	var user *LoginInfo
	user, err := ls.userRepo.GetInstance(userName)
	if err != nil {
		return nil, err
	}
	if user == nil {
		//add unknown user as a new user
		user = createNewUser(ls.hasher, userName, passwd)
		err = ls.Add(user)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	//user exists already, check password
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
