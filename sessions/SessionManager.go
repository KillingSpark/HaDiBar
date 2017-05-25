package sessions

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

//Session represents a Session
//The token originates from authenticating at the payment engine
type Session struct {
	Token string
	ID    string
	Name  string
}

//SessionManager handles the sessions
type SessionManager struct {
	sessionmap map[string]Session
	lock       sync.Mutex
}

//NewSessionManager does it
func NewSessionManager() *SessionManager {
	sm := SessionManager{}
	sm.sessionmap = make(map[string]Session)
	return &sm
}

//CheckSession checks if the token is valid and then executes the given handle
func (manager *SessionManager) CheckSession(ss *SessionManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID := ctx.Request.Header.Get("sessionID")

		if sessionID == "" {
			println("no session header found. Adding new one")
			ctx.Writer.Header().Set("sessionID", ss.MakeSessionID())
		} else {
			ctx.Writer.Header().Set("sessionID", sessionID)
			println("call from session: " + sessionID)
		}
		ctx.Writer.WriteHeader(http.StatusCreated)
		ctx.Next()
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

//MakeSessionID creates a new Session and returns the ID
func (manager *SessionManager) MakeSessionID() string {
	ID := "Lululululululu"
	session := Session{Token: "Blabla", ID: ID}
	manager.sessionmap[ID] = session
	return session.ID
}
