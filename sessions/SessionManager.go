package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/killingspark/HaDiBar/logger"
)

//Session represents a Session
//The token originates from authenticating at the payment engine
type Session struct {
	Token      string
	ID         string
	Name       string
	expiryDate int64
}

//SessionManager handles the sessions
type SessionManager struct {
	sessionmap  map[string]Session
	maxLifeTime int64
}

//NewSessionManager creates new SessionManager and starts GC as goroutine
func NewSessionManager() *SessionManager {
	sm := SessionManager{}
	sm.sessionmap = make(map[string]Session)
	sm.maxLifeTime = 24 * 60 * 60 * 1000 * 1000 * 1000
	go sm.GC()
	return &sm
}

//CheckSession checks if the token is valid and then executes the given handle
func (manager *SessionManager) CheckSession(ctx *gin.Context) {
	var sessionID = ctx.Request.Header.Get("sessionID")

	if sessionID == "" {
		logger.Logger.Debug("no session header found. Adding new one")
		sessionID = manager.MakeSessionID()
	} else {
		logger.Logger.Debug("call from session: " + sessionID)
	}

	//headers get written by gin
	ctx.Writer.Header().Set("sessionID", sessionID)
	session, err := manager.GetSession(sessionID)
	if err == nil {
		ctx.Set("session", session)
		ctx.Next()
	} else {
		logger.Logger.Warning(err.Error())
		ctx.Writer.WriteString("invalid session")
	}
}

//GetSession returns a session for the id or an error
func (manager *SessionManager) GetSession(id string) (Session, error) {
	val, ok := manager.sessionmap[id]
	if !ok {
		return Session{}, errors.New("session: " + id + " is not present")
	}
	return val, nil
}

//GC Starts the garbage collection for expired sessions
func (manager *SessionManager) GC() {
	for key, ses := range manager.sessionmap {
		if time.Now().UnixNano()-ses.expiryDate >= manager.maxLifeTime {
			delete(manager.sessionmap, key)
		}
	}
	time.AfterFunc(time.Duration(manager.maxLifeTime), manager.GC)
}

//MakeSessionID creates a new Session and returns the ID
func (manager *SessionManager) MakeSessionID() string {
	ID := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, ID); err != nil {
		return ""
	}
	encID := base64.URLEncoding.EncodeToString(ID)
	session := Session{Token: "", ID: encID, expiryDate: time.Now().UnixNano() - manager.maxLifeTime}
	manager.sessionmap[encID] = session
	return session.ID
}
