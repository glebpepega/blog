package server

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

type IndexHTML struct {
	Supn          SignUpNotification
	Sinn          SignInNotification
	SessionExists bool
	Login         string
}

type SignUpNotification struct {
	SupNExists bool
	SupNText   string
}

type SignInNotification struct {
	SinNExists bool
	SinNText   string
}

type BlogHTML struct {
	Cnts      []Cnt
	CntExists bool
	CntString string
	Login     string
}

type Cnt struct {
	Login       string
	Index       int
	TextIsLit   bool
	InsertImage bool
	Text        string
}

func (s *server) root(w http.ResponseWriter, r *http.Request) {
	iHTML := newIndexHTML()
	switch r.Method {
	case "GET":
		login, _, exists := validCookieExists(s, r)
		if exists {
			iHTML.SessionExists = true
			iHTML.Login = login
		} else {
			terminateSessionIDCookieIfExists(w, r)
		}
		if err := constructHTML("static/index.html", w, iHTML); err != nil {
			log.Println(err)
		}
	case "POST":
		if len(r.FormValue("password")) < 6 {
			iHTML.Supn.SupNExists = true
			iHTML.Supn.SupNText = "password must be at least 6 characters long"
			if err := constructHTML("static/index.html", w, iHTML); err != nil {
				log.Println(err)
			}
		} else if err := s.db.RegisterUser(r.FormValue("login"), r.FormValue("password")); err != nil {
			iHTML.Supn.SupNExists = true
			if err.Error() == "user already exists" {
				iHTML.Supn.SupNText = err.Error()
			} else {
				iHTML.Supn.SupNText = "username or password must be less than 50 characters long"
			}
			if err := constructHTML("static/index.html", w, iHTML); err != nil {
				log.Println(err)
			}
		} else {
			iHTML.Supn.SupNExists = true
			iHTML.Supn.SupNText = "account created!"
			if err := constructHTML("static/index.html", w, iHTML); err != nil {
				log.Println(err)
			}
		}
	}
}

func (s *server) logOut(w http.ResponseWriter, r *http.Request) {
	iHTML := newIndexHTML()
	terminateSessionIDCookieIfExists(w, r)
	if err := constructHTML("static/index.html", w, iHTML); err != nil {
		log.Println(err)
	}
}

func (s *server) fileServer() http.Handler {
	return http.FileServer(http.Dir("./static"))
}

func (s *server) myPage(w http.ResponseWriter, r *http.Request) {
	iHTML := newIndexHTML()
	bHtml := newBlogHTML()
	switch r.Method {
	case "GET":
		login, cnt, exists := validCookieExists(s, r)
		if exists {
			bHtml.Login = login
			bHtml.CntString = cnt
			constructBlog(bHtml, s, w, r)
		} else {
			terminateSessionIDCookieIfExists(w, r)
			if err := constructHTML("static/index.html", w, iHTML); err != nil {
				log.Println(err)
			}
		}
	case "POST":
		cnt, err := s.db.SearchUser(r.FormValue("login"), r.FormValue("password"))
		if err != nil {
			iHTML.Sinn.SinNExists = true
			iHTML.Sinn.SinNText = err.Error()
			if err := constructHTML("static/index.html", w, iHTML); err != nil {
				log.Println(err)
			}
		} else {
			login, cnt_, exists := validCookieExists(s, r)
			if exists {
				bHtml.Login = login
				bHtml.CntString = cnt_
				constructBlog(bHtml, s, w, r)
			} else {
				bHtml.Login = r.FormValue("login")
				bHtml.CntString = cnt
				newSessionLogIn(bHtml, s, w, r)
			}
		}
	}
}

func (s *server) newBlog(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	login, _, exists := validCookieExists(s, r)
	if exists {
		message := strings.ReplaceAll(r.FormValue("message"), "\n", " ")
		cnt, err := s.db.NewBlog(login, message)
		if err != nil {
			log.Println(err)
		}
		bHTML := newBlogHTML()
		bHTML.Login = login
		bHTML.CntString = cnt
		constructBlog(bHTML, s, w, r)
	} else {
		terminateSessionIDCookieIfExists(w, r)
		iHTML := newIndexHTML()
		if err := constructHTML("static/index.html", w, iHTML); err != nil {
			log.Println(err)
		}
	}
}

func (s *server) removeAPost(w http.ResponseWriter, r *http.Request) {
	login, _, exists := validCookieExists(s, r)
	if exists {
		index, err := strconv.Atoi(r.FormValue("index"))
		if err != nil {
			log.Println(err)
		}
		cnt, err := s.db.RemoveAPost(login, index)
		if err != nil {
			log.Println(err)
		}
		bHTML := newBlogHTML()
		bHTML.Login = login
		bHTML.CntString = cnt
		constructBlog(bHTML, s, w, r)
	} else {
		iHTML := newIndexHTML()
		terminateSessionIDCookieIfExists(w, r)
		if err := constructHTML("static/index.html", w, iHTML); err != nil {
			log.Println(err)
		}
	}
}

func (s *server) removePosts(w http.ResponseWriter, r *http.Request) {
	login, _, exists := validCookieExists(s, r)
	if exists {
		if err := s.db.RemovePosts(login); err != nil {
			log.Println(err)
		}
		bHTML := newBlogHTML()
		bHTML.Login = login
		bHTML.CntString = ""
		constructBlog(bHTML, s, w, r)
	} else {
		terminateSessionIDCookieIfExists(w, r)
		iHTML := newIndexHTML()
		if err := constructHTML("static/index.html", w, iHTML); err != nil {
			log.Println(err)
		}
	}
}
