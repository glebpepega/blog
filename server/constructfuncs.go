package server

import (
	"html/template"
	"io"
	"log"
	"net/http"
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
	Contents  []*Post
	CntsExist bool
	Login     string
	Dapn      DeleteAllPostsNotification
}

type Post struct {
	ID        string
	Author    string
	Timestamp string
	Text      string
	Images    []Image
}

type Image struct {
	Link string
}

type DeleteAllPostsNotification struct {
	DAPNExists bool
	DAPNText   string
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
	if len(bHTML.Contents) > 0 {
		bHTML.CntsExist = true
	}
	if err := constructHTML("static/blog.html", w, bHTML); err != nil {
		log.Println(err)
	}
}
