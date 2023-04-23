package s

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Notification struct {
	Supn SignUpNotification
	Sinn SignInNotification
}

type SignUpNotification struct {
	SupNExists bool
	SupNText   string
}

type SignInNotification struct {
	SinNExists bool
	SinNText   string
}

type BlogCnt struct {
	CntExists bool
	Cnts      []Cnt
	Login     string
	Password  string
}

type Cnt struct {
	Login       string
	Index       int
	TextIsLit   bool
	InsertImage bool
	Cnt         string
}

func (s *server) root(w http.ResponseWriter, r *http.Request) {
	nf := Notification{}
	switch r.Method {
	case "GET":
		if err := constructHTML("static/index.html", w, nf); err != nil {
			log.Println(err)
		}
	case "POST":
		if len(r.FormValue("password")) < 6 {
			nf.Supn.SupNExists = true
			nf.Supn.SupNText = "password must be at least 6 characters long"
			if err := constructHTML("static/index.html", w, nf); err != nil {
				log.Println(err)
			}
		} else if err := s.db.RegisterUser(r.FormValue("login"), r.FormValue("password")); err != nil {
			nf.Supn.SupNExists = true
			if err.Error() == "user already exists" {
				nf.Supn.SupNText = err.Error()
			} else {
				nf.Supn.SupNText = "username or password must be less than 50 characters long"
			}
			if err := constructHTML("static/index.html", w, nf); err != nil {
				log.Println(err)
			}
		} else {
			nf.Supn.SupNExists = true
			nf.Supn.SupNText = "account created!"
			if err := constructHTML("static/index.html", w, nf); err != nil {
				log.Println(err)
			}
		}
	}
}

func (s *server) fileServer() http.Handler {
	return http.FileServer(http.Dir("./static"))
}

func (s *server) myPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	nf := Notification{}
	cnt, err := s.db.SearchUser(r.FormValue("login"), r.FormValue("password"))
	if err != nil {
		nf.Sinn.SinNExists = true
		nf.Sinn.SinNText = err.Error()
		if err := constructHTML("static/index.html", w, nf); err != nil {
			log.Println(err)
		}
	} else {
		constructBlog(cnt, s, w, r)
	}
}

func (s *server) newBlog(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	message := strings.ReplaceAll(r.FormValue("message"), "\n", " ")
	cnt, err := s.db.NewBlog(r.FormValue("login"), message)
	if err != nil {
		log.Println(err)
	}
	constructBlog(cnt, s, w, r)
}

func (s *server) removePosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err := s.db.RemovePosts(r.FormValue("login")); err != nil {
		log.Println(err)
	}
	constructBlog("", s, w, r)
}

func (s *server) removeAPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	index, err := strconv.Atoi(r.FormValue("index"))
	if err != nil {
		log.Println(err)
	}
	cnt, err := s.db.RemoveAPost(index, r.FormValue("login"))
	if err != nil {
		log.Println(err)
	}
	constructBlog(cnt, s, w, r)
}

func constructHTML(filename string, w io.Writer, data any) error {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(w, data); err != nil {
		return err
	}
	return nil
}

func constructBlog(cnt string, s *server, w http.ResponseWriter, r *http.Request) {
	bc := BlogCnt{}
	bc.Login = r.FormValue("login")
	bc.Password = r.FormValue("password")
	if cnt != "" {
		bc.CntExists = true
		scnSlice := strings.Split(cnt, "\n")
		for i, v := range scnSlice {
			c := Cnt{Cnt: v}
			c.Login = r.FormValue("login")
			c.Index = i
			if i%2 == 0 {
				c.TextIsLit = true
			}
			if i%2 == 1 {
				c.InsertImage = true
			}
			bc.Cnts = append(bc.Cnts, c)
			c.TextIsLit = false
			c.InsertImage = false
		}
	}
	if err := constructHTML("static/blog.html", w, bc); err != nil {
		log.Println(err)
	}
}
