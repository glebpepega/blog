package s

import (
	"log"
	"net/http"

	"github.com/glebpepega/blog/db"
)

type server struct {
	db *db.DB
	r  *http.ServeMux
}

func New() *server {
	return &server{
		db: db.NewDB(),
		r:  http.NewServeMux(),
	}
}

func (s *server) Start() {
	s.db.Start()
	s.r.Handle("/", http.HandlerFunc(s.root))
	s.r.Handle("/images/", http.StripPrefix("/images", s.fileServer()))
	s.r.Handle("/mypage/images/", http.StripPrefix("/mypage/images", s.fileServer()))
	s.r.HandleFunc("/mypage", s.myPage)
	s.r.HandleFunc("/mypage/newblog", s.newBlog)
	s.r.HandleFunc("/mypage/removeposts", s.removePosts)
	s.r.HandleFunc("/mypage/removeapost", s.removeAPost)
	log.Fatal(http.ListenAndServe(":8080", s.r))
}
