package services

import (
	"errors"
)

//Session represents a Session
//The token originates from authenticating at the payment engine
type Session struct {
	Token string
	ID    string
	Name  string
}

//SessionService handles the sessions
type SessionService struct {
	sessionmap map[string]Session
}

//MakeSessionService does it
func MakeSessionService() SessionService {
	ss := SessionService{}
	ss.sessionmap = make(map[string]Session)
	return ss
}

//GetSession returns a session for the id or an error
func (service *SessionService) GetSession(id string) (Session, error) {
	val, ok := service.sessionmap[id]
	if !ok {
		return Session{}, errors.New("session: " + id + " is not present")
	}
	return val, nil
}

//MakeSessionID creates a new Session and returns the ID
func (service *SessionService) MakeSessionID() string {
	ID := "Lululululululu"
	session := Session{Token: "Blabla", ID: ID}
	service.sessionmap[ID] = session
	return session.ID
}
