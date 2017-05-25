package accounts

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"time"
)

type token struct {
	value      string
	expiredate int64
}

//creates new token with an expiredate of now + 24h
func makeToken(value string) token {
	tok := token{value: value, expiredate: time.Now().UnixNano() - 24*60*60*1000*1000*1000}
	return tok
}

//LoginService handles all operations connected to identification
type LoginService struct {
	tokenmap map[string]Entity
	tokens   []token
}

//GetEntityFromToken returns the entity that belongs to the token. If the token is invalid/expired the boolean
//is going to be false
func (service *LoginService) GetEntityFromToken(tokenval string) (Entity, bool) {

	_, ok := service.lookUpToken(tokenval)
	if !ok {
		return Entity{}, false
	}

	ent := service.tokenmap[tokenval]

	if &ent == nil {

	}

	return ent, true
}

//lookup if the token is known to the server
func (service *LoginService) lookUpToken(tokenval string) (token, bool) {
	for index, tk := range service.tokens {
		if tk.expiredate >= time.Now().UnixNano() {
			//delete expired tokens
			println("Expired: " + tk.value)
			service.tokens = append(service.tokens[:index], service.tokens[index+1:]...)
		} else {
			if tk.value == tokenval {
				return tk, true
			}
		}
	}
	return token{}, false
}

//IsTokenValid checks if the token is valid
func (service *LoginService) IsTokenValid(tokenval string) bool {
	_, isValid := service.lookUpToken(tokenval)
	return isValid
}

//RequestToken returns a token for the credentials
func (service *LoginService) RequestToken(name, password string) (string, bool) {
	tokenstring := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, tokenstring); err != nil {
		return "", false
	}
	enctokenstring := base64.URLEncoding.EncodeToString(tokenstring)
	service.tokens = append(service.tokens, makeToken(enctokenstring))

	return enctokenstring, true
}
