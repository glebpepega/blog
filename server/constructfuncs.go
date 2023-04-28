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
	AlreadyAddedToFavPage bool
	Login                 string
	CreatedOn             string
	FavouritePages        []FavPage
	FavPsExist            bool
	Contents              []Post
	CntsExist             bool
	Dapn                  DeleteAllPostsNotification
}

type FavPage struct {
	Username string
}

type Post struct {
	ID       string
	Author   string
	PostedOn string
	Text     string
	Images   []Image
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

func constructBlog(filename string, data any, s *server, w http.ResponseWriter, r *http.Request) {
	bHTML := data.(*BlogHTML)
	if len(bHTML.FavouritePages) > 0 {
		bHTML.FavPsExist = true
	}
	if len(bHTML.Contents) > 0 {
		bHTML.CntsExist = true
	}
	if err := constructHTML(filename, w, bHTML); err != nil {
		log.Println(err)
	}
}
