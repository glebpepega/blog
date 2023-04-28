package server

import (
	"log"
	"net/http"
	"strings"
)

func (s *server) root(w http.ResponseWriter, r *http.Request) {
	iHTML := newIndexHTML()
	switch r.Method {
	case "GET":
		sID, exists := validSessionExists(s, r)
		if exists {
			login, err := s.db.getHomePageContents(sID)
			if err != nil {
				log.Println(err)
			}
			iHTML.SessionExists = true
			iHTML.Login = login
		} else {
			terminateSessionIDCookieIfExists(w, r)
		}
		if err = constructHTML("static/index.html", w, iHTML); err != nil {
			log.Println(err)
		}
	case "POST":
		if len(r.FormValue("password")) < 6 {
			iHTML.Supn.SupNExists = true
			iHTML.Supn.SupNText = "password must be at least 6 characters long"
		} else if err = s.db.registerUser(r.FormValue("login"), r.FormValue("password")); err != nil {
			iHTML.Supn.SupNExists = true
			if err.Error() == "user already exists" {
				iHTML.Supn.SupNText = err.Error()
			} else {
				iHTML.Supn.SupNText = "username or password must be less than 50 characters long"
			}
		} else {
			iHTML.Supn.SupNExists = true
			iHTML.Supn.SupNText = "account created!"
		}
		if err = constructHTML("static/index.html", w, iHTML); err != nil {
			log.Println(err)
		}
	}
}

func (s *server) logOut(w http.ResponseWriter, r *http.Request) {
	terminateSessionIDCookieIfExists(w, r)
	iHTML := newIndexHTML()
	if err = constructHTML("static/index.html", w, iHTML); err != nil {
		log.Println(err)
	}
}

func (s *server) fileServer() http.Handler {
	return http.FileServer(http.Dir("./static"))
}

func (s *server) myPage(w http.ResponseWriter, r *http.Request) {
	iHTML := newIndexHTML()
	bHTML := newBlogHTML()
	switch r.Method {
	case "GET":
		sID, exists := validSessionExists(s, r)
		if exists {
			u, err := s.db.getBlogContents(sID)
			if err != nil {
				log.Println(err)
			}
			bHTML.Login = u.Login
			bHTML.CreatedOn = u.CreatedOn
			bHTML.FavouritePages = u.FavouritePages
			bHTML.Contents = u.Blog
			constructBlog("static/blog.html", bHTML, s, w, r)
		} else {
			terminateSessionIDCookieIfExists(w, r)
			if err = constructHTML("static/index.html", w, iHTML); err != nil {
				log.Println(err)
			}
		}
	case "POST":
		if err = s.db.searchUser(r.FormValue("login"), r.FormValue("password")); err != nil {
			iHTML.Sinn.SinNExists = true
			iHTML.Sinn.SinNText = err.Error()
			if err = constructHTML("static/index.html", w, iHTML); err != nil {
				log.Println(err)
			}
		} else {
			newSessionID := generateSessionID()
			setSessionIDCookie(newSessionID, w, r)
			if err = s.db.addSession(r.FormValue("login"), newSessionID); err != nil {
				log.Println(err)
			}
			u, err := s.db.getBlogContents(newSessionID)
			if err != nil {
				log.Println(err)
			}
			bHTML.Login = u.Login
			bHTML.CreatedOn = u.CreatedOn
			bHTML.FavouritePages = u.FavouritePages
			bHTML.Contents = u.Blog
			constructBlog("static/blog.html", bHTML, s, w, r)
		}
	}
}

func (s *server) newBlog(w http.ResponseWriter, r *http.Request) {
	sID, exists := validSessionExists(s, r)
	if exists {
		imageLinks := strings.Split(r.FormValue("images"), "\n")
		if err := s.db.addNewPost(imageLinks, r.FormValue("message"), sID); err != nil {
			log.Println(err)
		}
		u, err := s.db.getBlogContents(sID)
		if err != nil {
			log.Println(err)
		}
		bHTML := newBlogHTML()
		bHTML.Login = u.Login
		bHTML.CreatedOn = u.CreatedOn
		bHTML.FavouritePages = u.FavouritePages
		bHTML.Contents = u.Blog
		constructBlog("static/blog.html", bHTML, s, w, r)
	} else {
		terminateSessionIDCookieIfExists(w, r)
		iHTML := newIndexHTML()
		if err = constructHTML("static/index.html", w, iHTML); err != nil {
			log.Println(err)
		}
	}
}

func (s *server) removeAPost(w http.ResponseWriter, r *http.Request) {
	sID, exists := validSessionExists(s, r)
	if exists {
		if err = s.db.removeAPost(r.FormValue("postid")); err != nil {
			log.Println(err)
		}
		u, err := s.db.getBlogContents(sID)
		if err != nil {
			log.Println(err)
		}
		bHTML := newBlogHTML()
		bHTML.Login = u.Login
		bHTML.CreatedOn = u.CreatedOn
		bHTML.FavouritePages = u.FavouritePages
		bHTML.Contents = u.Blog
		constructBlog("static/blog.html", bHTML, s, w, r)
	} else {
		terminateSessionIDCookieIfExists(w, r)
		iHTML := newIndexHTML()
		if err = constructHTML("static/index.html", w, iHTML); err != nil {
			log.Println(err)
		}
	}
}

func (s *server) removePosts(w http.ResponseWriter, r *http.Request) {
	sID, exists := validSessionExists(s, r)
	if exists {
		bHTML := newBlogHTML()
		if len(r.FormValue("password")) < 6 {
			bHTML.Dapn.DAPNExists = true
			bHTML.Dapn.DAPNText = "password must be at least 6 characters long"
		} else if err = s.db.removePosts(r.FormValue("password"), sID); err != nil {
			if err.Error() == "wrong password" {
				bHTML.Dapn.DAPNExists = true
				bHTML.Dapn.DAPNText = err.Error()
			} else {
				log.Println(err)
			}
		} else {
			bHTML.Dapn.DAPNExists = true
			bHTML.Dapn.DAPNText = "posts deleted"
		}
		u, err := s.db.getBlogContents(sID)
		if err != nil {
			log.Println(err)
		}
		bHTML.Login = u.Login
		bHTML.CreatedOn = u.CreatedOn
		bHTML.FavouritePages = u.FavouritePages
		bHTML.Contents = u.Blog
		constructBlog("static/blog.html", bHTML, s, w, r)
	} else {
		terminateSessionIDCookieIfExists(w, r)
		iHTML := newIndexHTML()
		if err = constructHTML("static/index.html", w, iHTML); err != nil {
			log.Println(err)
		}
	}
}

func (s *server) visit(w http.ResponseWriter, r *http.Request) {
	sID, exists := validSessionExists(s, r)
	if exists {
		bHTML := newBlogHTML()
		bHTML.Login = r.FormValue("visitlogin")
		if err = s.db.searchUser(bHTML.Login, ""); err.Error() != "no such user" {
			if err = s.db.isPageInFavourites(bHTML.Login, sID); err == nil {
				bHTML.AlreadyAddedToFavPage = true
			}
			u, err := s.db.getVisitContents(bHTML.Login)
			if err != nil {
				log.Println(err)
			}
			bHTML.CreatedOn = u.CreatedOn
			bHTML.Contents = u.Blog
			constructBlog("static/visit.html", bHTML, s, w, r)
		} else {
			u, err := s.db.getBlogContents(sID)
			if err != nil {
				log.Println(err)
			}
			bHTML.Login = u.Login
			bHTML.CreatedOn = u.CreatedOn
			bHTML.FavouritePages = u.FavouritePages
			bHTML.Contents = u.Blog
			constructBlog("static/blog.html", bHTML, s, w, r)
		}
	} else {
		terminateSessionIDCookieIfExists(w, r)
		iHTML := newIndexHTML()
		if err = constructHTML("static/index.html", w, iHTML); err != nil {
			log.Println(err)
		}
	}
}

func (s *server) addToFavs(w http.ResponseWriter, r *http.Request) {
	sID, exists := validSessionExists(s, r)
	if exists {
		bHTML := newBlogHTML()
		bHTML.Login = r.FormValue("visitlogin")
		if err = s.db.addPageToFavourites(bHTML.Login, sID); err != nil {
			log.Println(err)
		} else {
			bHTML.AlreadyAddedToFavPage = true
		}
		u, err := s.db.getVisitContents(bHTML.Login)
		if err != nil {
			log.Println(err)
		}
		bHTML.CreatedOn = u.CreatedOn
		bHTML.Contents = u.Blog
		constructBlog("static/visit.html", bHTML, s, w, r)
	} else {
		terminateSessionIDCookieIfExists(w, r)
		iHTML := newIndexHTML()
		if err = constructHTML("static/index.html", w, iHTML); err != nil {
			log.Println(err)
		}
	}
}

func (s *server) removeFromFavs(w http.ResponseWriter, r *http.Request) {
	sID, exists := validSessionExists(s, r)
	if exists {
		bHTML := newBlogHTML()
		bHTML.Login = r.FormValue("visitlogin")
		if err = s.db.removePageFromFavourites(bHTML.Login, sID); err != nil {
			log.Println(err)
		} else {
			bHTML.AlreadyAddedToFavPage = false
		}
		u, err := s.db.getVisitContents(bHTML.Login)
		if err != nil {
			log.Println(err)
		}
		bHTML.CreatedOn = u.CreatedOn
		bHTML.Contents = u.Blog
		constructBlog("static/visit.html", bHTML, s, w, r)
	} else {
		terminateSessionIDCookieIfExists(w, r)
		iHTML := newIndexHTML()
		if err = constructHTML("static/index.html", w, iHTML); err != nil {
			log.Println(err)
		}
	}
}
