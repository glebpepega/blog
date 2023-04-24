package server

import (
	"log"
	"net/http"

	"github.com/glebpepega/blog/database"
)

type server struct {
	db *database.DB
	r  *http.ServeMux
}

func New() *server {
	return &server{
		db: database.NewDB(),
		r:  http.NewServeMux(),
	}
}

func (s *server) Start() {
	s.db.Start()
	s.r.Handle("/", http.HandlerFunc(s.root))
	s.r.HandleFunc("/logout", s.logOut)
	s.r.Handle("/images/", http.StripPrefix("/images", s.fileServer()))
	s.r.Handle("/mypage/images/", http.StripPrefix("/mypage/images", s.fileServer()))
	s.r.HandleFunc("/mypage", s.myPage)
	s.r.HandleFunc("/mypage/newblog", s.newBlog)
	s.r.HandleFunc("/mypage/removeapost", s.removeAPost)
	s.r.HandleFunc("/mypage/removeposts", s.removePosts)
	log.Fatal(http.ListenAndServe(":8080", s.r))
}