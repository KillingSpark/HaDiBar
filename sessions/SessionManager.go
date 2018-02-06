package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/killingspark/HaDiBar/logger"
	"github.com/killingspark/HaDiBar/restapi"
)

//Session represents a Session
//The token originates from authenticating at the payment engine
type Session struct {
	Token      string
	ID         string
	Name       string
	Floor      string
	expiryDate int64
}

//SessionManager handles the sessions
type SessionManager struct {
	sessionmap  map[string]*Session
	maxLifeTime int64
}

//NewSessionManager creates new SessionManager and starts GC as goroutine
func NewSessionManager() *SessionManager {
	sm := SessionManager{}
	sm.sessionmap = make(map[string]*Session)
	sm.maxLifeTime = 24 * 60 * 60 * 1000 * 1000 * 1000
	go sm.GC()
	return &sm
}

//IsSessionLoggedIn checks wether the given session should be seen as logged in
func (manager *SessionManager) IsSessionLoggedIn(sessionID string) bool {
	session, err := manager.GetSession(sessionID)

	if err != nil {
		logger.Logger.Debug(sessionID + " tried access but isnt valid")
		return false
	}

	if session.Token != "" && session.Name != "" {
		return true
	}

	logger.Logger.Debug(sessionID + " tried to login but isnt logged in")
	return false
}

//CheckLoginStatus checks if the session is logged in and then executes the given handle
func (manager *SessionManager) CheckLoginStatus(ctx *gin.Context) {
	var sessionID = ctx.Request.Header.Get("sessionID")
	if manager.IsSessionLoggedIn(sessionID) {
		logger.Logger.Debug("Logincheck good for: " + sessionID)
		ctx.Next()
	} else {
		logger.Logger.Warning("Logincheck bad for: " + sessionID)
		response, _ := restapi.NewErrorResponse("You must be logged in here").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
	}
}

//CheckSession checks if the token is valid and then executes the given handle
func (manager *SessionManager) CheckSession(ctx *gin.Context) {
	var sessionID = ctx.Request.Header.Get("sessionID")

	if sessionID == "" {
		logger.Logger.Debug("no session header found. Adding new one")
		sessionID = manager.MakeSessionID()
	} else {
		logger.Logger.Debug("call from session: " + sessionID + " to URL: " + ctx.Request.URL.RawPath)
	}

	//headers get written by gin
	ctx.Writer.Header().Set("sessionID", sessionID)
	session, err := manager.GetSession(sessionID)
	if err == nil {
		ctx.Set("session", session)
		ctx.Next()
	} else {
		logger.Logger.Warning(err.Error())
		response, _ := restapi.NewErrorResponse("No valid session").Marshal()
		fmt.Fprint(ctx.Writer, string(response))
		ctx.Abort()
	}
}

//SetSessionToken set the token for the session
func (manager *SessionManager) SetSessionToken(id string, token string) {
	val, ok := manager.sessionmap[id]
	if ok {
		val.Token = token
	}
}

//SetSessionFloor set the token for the session
func (manager *SessionManager) SetSessionFloor(id string, floor string) {
	val, ok := manager.sessionmap[id]
	if ok {
		val.Floor = floor
	}
}

//SetSessionName sets the name for the session login
func (manager *SessionManager) SetSessionName(id string, name string) {
	val, ok := manager.sessionmap[id]
	if ok {
		val.Name = name
	}
}

//GetSession returns a session for the id or an error
func (manager *SessionManager) GetSession(id string) (*Session, error) {
	val, ok := manager.sessionmap[id]
	if !ok {
		return &Session{}, errors.New("session: " + id + " is not present")
	}
	return val, nil
}

//GC Starts the garbage collection for expired sessions
func (manager *SessionManager) GC() {
	for key, ses := range manager.sessionmap {
		if manager.isExspired(ses.expiryDate) {
			delete(manager.sessionmap, key)
		}
	}
	time.AfterFunc(time.Duration(manager.maxLifeTime), manager.GC)
}

func (manager *SessionManager) isExspired(timestamp int64) bool {
	return time.Now().UnixNano()-timestamp >= manager.maxLifeTime
}

//MakeSessionID creates a new Session and returns the ID
func (manager *SessionManager) MakeSessionID() string {
	ID := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, ID); err != nil {
		return ""
	}
	encID := base64.URLEncoding.EncodeToString(ID)
	session := Session{Token: "", ID: encID, expiryDate: time.Now().UnixNano() - manager.maxLifeTime}
	manager.sessionmap[encID] = &session
	return session.ID
}
