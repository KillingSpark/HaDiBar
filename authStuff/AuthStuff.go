package authStuff

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/killingspark/HaDiBar/logger"
	"github.com/killingspark/HaDiBar/restapi"
	"github.com/killingspark/HaDiBar/settings"
)

//Entity (s) represent owners of an Account
type LoginInfo struct {
	Name     string
	LoggedIn bool
	Salt     string
	Pwhash   string
}

type Authentikator interface {
	isValid(id, pw string) (*LoginInfo, error)
}

type Session struct {
	id   string
	info *LoginInfo
}

type Auth struct {
	sessionMap map[string](*Session)
	tester     Authentikator
}

func NewAuth() (*Auth, error) {
	auth := &Auth{}
	auth.sessionMap = make(map[string](*Session))
	var err error
	auth.tester, err = NewLoginService(settings.S.DataDir)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func (auth *Auth) AddNewSession() string {
	ID := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, ID); err != nil {
		return ""
	}
	encID := base64.URLEncoding.EncodeToString(ID)

	session := Session{id: encID}
	auth.sessionMap[encID] = &session
	return session.id
}

var ErrAlreadyLoggedIn = errors.New("Already logged in")

func (auth *Auth) LogIn(id, name, password string) error {
	session, ok := auth.sessionMap[id]
	if !ok {
		return ErrInvalidSession
	}

	if session.info != nil && session.info.LoggedIn {
		return ErrAlreadyLoggedIn
	}

	newinfo, err := auth.tester.isValid(name, password)

	if err != nil {
		return err
	}

	session.info = newinfo
	session.info.LoggedIn = true
	session.info.Name = name
	return nil
}

func (auth *Auth) GetSessionInfo(id string) (*LoginInfo, error) {
	session, err := auth.getSession(id)
	if err != nil {
		return nil, err
	}
	return session.info, nil
}

var ErrInvalidSession = errors.New("Session not valid")

func (auth *Auth) LogOut(id string) error {
	session, ok := auth.sessionMap[id]
	if !ok {
		return ErrInvalidSession
	}

	session.info = &LoginInfo{LoggedIn: false}
	return nil
}

//CheckSession checks if the token is valid and then executes the given handle
func (auth *Auth) CheckSession(ctx *gin.Context) {
	var sessionID = ctx.Request.Header.Get("sessionID")

	if sessionID == "" {
		logger.Logger.Debug("no session header found. Adding new one")
		sessionID = auth.AddNewSession()
	} else {
		logger.Logger.Debug("call from session: " + sessionID + " to URL: " + ctx.Request.URL.RawPath)
	}

	//headers get written by gin
	ctx.Writer.Header().Set("sessionID", sessionID)
	session, err := auth.getSession(sessionID)
	if err == nil {
		ctx.Set("session", session)
		ctx.Next()
	} else {
		logger.Logger.Warning(err.Error())
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
		logger.Logger.Debug("Logincheck good for: " + sessionID)
		ctx.Set("logininfo", session.info)
		ctx.Next()
	} else {
		logger.Logger.Warning("Logincheck bad for: " + sessionID)
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
