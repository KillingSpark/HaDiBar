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
	user, _ := ls.userRepo.GetInstance(new.Name)
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
	user.Salt = SaltPw(hasher, strconv.FormatInt(time.Now().UnixNano()%rand.Int63(), 10), username)
	user.Pwhash = SaltPw(hasher, passwd, user.Salt)
	return user
}

func (ls *LoginService) isValid(userName, passwd string) (*LoginInfo, error) {
	var user *LoginInfo
	user, err := ls.userRepo.GetInstance(userName)
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
	if SaltPw(ls.hasher, passwd, user.Salt) == user.Pwhash {
		user.LastLogin = time.Now()
		ls.userRepo.SaveInstance(user)
		return user, nil
	}
	return nil, ErrWrongCredetials
}

func SaltPw(hasher hash.Hash, pw, salt string) string {
	hasher.Reset()
	hasher.Write([]byte(pw + salt))
	saltedpw := make([]byte, 4*int(math.Ceil((float64(32)/3))))

	base64.StdEncoding.Encode(saltedpw, hasher.Sum(nil))
	return string(saltedpw)
}
