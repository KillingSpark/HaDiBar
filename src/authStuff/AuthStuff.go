package authStuff

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/killingspark/hadibar/src/restapi"
)

//LoginInfo is passed into the context if the session has a logged in user
type LoginInfo struct {
	Name      string
	LoggedIn  bool
	Salt      string
	Pwhash    string
	LastLogin time.Time
	Email     string
}

//Authentikator is an interface that will allow for other sign-in methods later
type Authentikator interface {
	isValid(id, pw string) (*LoginInfo, error)
}

//Session identifies a session with a Client. If the client logs in, the session remembers the login info (without the password of course) until the client logs out.
type Session struct {
	id         string
	info       *LoginInfo
	lastAction time.Time
}

//Auth maps session ids to sessions
type Auth struct {
	sessionMap map[string](*Session)
	sessionTTL time.Duration
	ls         *LoginService
}

//NewAuth is a constructor for Auth
func NewAuth(datadir string, sessionTTL int) (*Auth, error) {
	auth := &Auth{}
	auth.sessionTTL = time.Duration(sessionTTL) * time.Second
	auth.sessionMap = make(map[string](*Session))
	var err error
	auth.ls, err = NewLoginService(datadir)
	if err != nil {
		return nil, err
	}
	go auth.cleanSessions()
	return auth, nil
}

func (auth *Auth) cleanSessions() {
	if auth.sessionTTL <= 0 {
		return
	}
	for {
		time.Sleep(1 * time.Minute)
		for key, ses := range auth.sessionMap {
			toRemove := make([]string, 0)
			if ses.lastAction.Add(auth.sessionTTL).Before(time.Now()) {
				toRemove = append(toRemove, key)
			}
			for _, key := range toRemove {
				delete(auth.sessionMap, key)
			}
		}
	}
}

//AddNewSession creates a new sessionid and remembers the session for later reference by the client
func (auth *Auth) AddNewSession() string {
	ID := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, ID); err != nil {
		return ""
	}
	encID := base64.URLEncoding.EncodeToString(ID)

	session := Session{id: encID, lastAction: time.Now()}
	auth.sessionMap[encID] = &session
	return session.id
}

var ErrAlreadyLoggedIn = errors.New("Already logged in")

//LogIn checks the credentials against the authentikator and marks the session as loggedin
func (auth *Auth) LogIn(sesID, name, password string) error {
	session, ok := auth.sessionMap[sesID]
	if !ok {
		return ErrInvalidSession
	}

	if session.info != nil && session.info.LoggedIn {
		return ErrAlreadyLoggedIn
	}

	newinfo, err := auth.ls.isValid(name, password)

	if err != nil {
		return err
	}

	session.info = newinfo
	session.info.LoggedIn = true
	session.info.Name = name
	return nil
}

//GetSessionInfo maps the sessionid to the LoginInfo
func (auth *Auth) GetSessionInfo(id string) (*LoginInfo, error) {
	session, err := auth.getSession(id)
	if err != nil {
		return nil, err
	}
	return session.info, nil
}

var ErrInvalidSession = errors.New("Session not valid")

//LogOut clears the LoginInfo of this session
func (auth *Auth) LogOut(id string) error {
	session, ok := auth.sessionMap[id]
	if !ok {
		return ErrInvalidSession
	}

	session.info = &LoginInfo{LoggedIn: false}
	return nil
}

//CheckSession checks if the sessionID is valid. If no sessionID is given, a new one is created and added as a header
func (auth *Auth) CheckSession(ctx *gin.Context) {
	var sessionID = ctx.Request.Header.Get("sessionID")

	if sessionID == "" {
		log.WithFields(log.Fields{}).Debug("No session header found. Adding new one")
		sessionID = auth.AddNewSession()
	} else {
		log.WithFields(log.Fields{"session": sessionID, "URL": ctx.Request.URL.String()}).Debug("Checked sessionid")
	}

	//headers get written by gin
	ctx.Writer.Header().Set("sessionID", sessionID)
	session, err := auth.getSession(sessionID)
	if err == nil {
		session.lastAction = time.Now()
		ctx.Set("session", session)
		ctx.Next()
	} else {
		log.WithFields(log.Fields{"session": sessionID, "Error": err.Error()}).Warn("Search SessionID error")
		response, _ := restapi.NewErrorResponse("No valid session").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}
}

//CheckLoginStatus checks if the session is logged in and then executes the given handle
func (auth *Auth) CheckLoginStatus(ctx *gin.Context) {
	sessionID := ctx.Request.Header.Get("sessionID")
	session, err := auth.getSession(sessionID)
	if err != nil {
		response, _ := restapi.NewNosesResponse("Sessionid invalid").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}
	if session.info != nil && session.info.LoggedIn {
		log.WithFields(log.Fields{"session": sessionID}).Debug("Login check good")
		ctx.Set("logininfo", session.info)
		ctx.Next()
	} else {
		log.WithFields(log.Fields{"session": sessionID, "URL": ctx.Request.URL.String()}).Warn("Login check bad")
		response, _ := restapi.NewErrorResponse("You must be logged in here").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
		return
	}
}

func (auth *Auth) getSession(id string) (*Session, error) {
	session, ok := auth.sessionMap[id]
	if ok {
		return session, nil
	}
	return nil, ErrInvalidSession
}

//GetLoginInfoFromCtx : Utility function for other controllers to get the LoginInfo from their Context
func GetLoginInfoFromCtx(ctx *gin.Context) (*LoginInfo, error) {
	var info *LoginInfo

	if inter, ok := ctx.Get("logininfo"); ok {
		info, ok = inter.(*LoginInfo)
		if !ok {
			return nil, errors.New("Not a LoginInfo while expecting LoginInfo. This is an internal misbehaviour. Contact an admin about this")
		}
	} else {
		return nil, errors.New("No Login-Info found. Try to log in again")
	}
	return info, nil
}
