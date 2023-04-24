package server

import (
	"html/template"
	"io"
	"log"
	"net/http"
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

func newIndexHTML() *IndexHTML {
	return &IndexHTML{}
}

func newBlogHTML() *BlogHTML {
	return &BlogHTML{}
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

func constructBlog(data any, s *server, w http.ResponseWriter, r *http.Request) {
	bHTML := data.(*BlogHTML)
	if bHTML.CntString != "" {
		bHTML.CntExists = true
		csnSlice := strings.Split(bHTML.CntString, "\n")
		for i, v := range csnSlice {
			cnt := Cnt{}
			cnt.Text = v
			cnt.Login = bHTML.Login
			cnt.Index = i
			if i%2 == 0 {
				cnt.TextIsLit = true
			}
			if i%2 == 1 {
				cnt.InsertImage = true
			}
			bHTML.Cnts = append(bHTML.Cnts, cnt)
			cnt.TextIsLit = false
			cnt.InsertImage = false
		}
	}
	if err := constructHTML("static/blog.html", w, bHTML); err != nil {
		log.Println(err)
	}
}

func newSessionLogIn(data any, s *server, w http.ResponseWriter, r *http.Request) {
	newSessionID := generateSessionID()
	setSessionIDCookie(newSessionID, w, r)
	s.db.AddSession(r.FormValue("login"), newSessionID)
	bHTML := data.(*BlogHTML)
	constructBlog(bHTML, s, w, r)
}
