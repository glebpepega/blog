package server

import (
	"math/rand"
	"net/http"
)

func setSessionIDCookie(newSessionID string, w http.ResponseWriter, r *http.Request) {
	sessionID := &http.Cookie{
		Name:   "SESSIONID",
		Value:  newSessionID,
		MaxAge: 86400,
	}
	http.SetCookie(w, sessionID)
}

func validSessionExists(s *server, r *http.Request) (sessionID string, status bool) {
	cookie, err := r.Cookie("SESSIONID")
	if err != nil {
		return "", false
	} else {
		if err := s.db.searchSession(cookie.Value); err != nil {
			return "", false
		} else {
			return cookie.Value, true
		}
	}
}

func terminateSessionIDCookieIfExists(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SESSIONID")
	if err == nil {
		cookie.MaxAge = -1
	}
	http.SetCookie(w, cookie)
}

func generateSessionID() (sessionID string) {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := 0; i < 10; i++ {
		sessionID += string(chars[rand.Intn(62)])
	}
	return sessionID
}
