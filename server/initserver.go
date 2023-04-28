package server

import (
	"log"
	"net/http"
)

type server struct {
	db *DB
	r  *http.ServeMux
}

func New() *server {
	return &server{
		db: NewDB(),
		r:  http.NewServeMux(),
	}
}

func (s *server) Start() {
	s.db.Start()
	s.r.HandleFunc("/", s.root)
	s.r.HandleFunc("/logout", s.logOut)
	s.r.Handle("/images/", http.StripPrefix("/images", s.fileServer()))
	s.r.Handle("/mypage/images/", http.StripPrefix("/mypage/images", s.fileServer()))
	s.r.HandleFunc("/mypage", s.myPage)
	s.r.HandleFunc("/mypage/newblog", s.newBlog)
	s.r.HandleFunc("/mypage/removeapost", s.removeAPost)
	s.r.HandleFunc("/mypage/removeposts", s.removePosts)
	log.Fatal(http.ListenAndServe(":8080", s.r))
}
