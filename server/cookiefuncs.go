package server

import (
	"math/rand"
	"net/http"
	"strconv"
)

func setSessionIDCookie(newSessionID string, w http.ResponseWriter, r *http.Request) {
	sessionID := &http.Cookie{
		Name:   "SESSIONID",
		Value:  newSessionID,
		MaxAge: 86400,
	}
	http.SetCookie(w, sessionID)
}

func validCookieExists(s *server, r *http.Request) (login string, content string, status bool) {
	cookie, err := r.Cookie("SESSIONID")
	if err != nil {
		return "", "", false
	} else {
		login, cnt, err := s.db.SearchSession(cookie.Value)
		return login, cnt, err == nil
	}
}

func terminateSessionIDCookieIfExists(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SESSIONID")
	if err == nil {
		cookie.MaxAge = -1
	}
	http.SetCookie(w, cookie)
}

func generateSessionID() string {
	sessionID := ""
	for i := 0; i < 10; i++ {
		sessionID += strconv.Itoa(rand.Intn(10))
	}
	return sessionID
}
