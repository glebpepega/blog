package server

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	ID             string
	Login          string
	Password       string
	CreatedOn      string
	FavouritePages []FavPage
	Blog           []Post
}

func newUser() User {
	return User{}
}

func newPost() Post {
	return Post{}
}

func (d *DB) registerUser(login, password string) error {
	createdOn := time.Now().Format("~~ 15:04:05 ~~ 02.01.2006 ~~")
	if _, err = d.db.Exec("INSERT INTO users (login, password,created_on) VALUES ($1, $2, $3);", login, password, createdOn); err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_pkey\"" {
			return fmt.Errorf("user already exists")
		}
		return err
	}
	return nil
}

func (d *DB) searchUser(login, password string) error {
	u := newUser()
	if err = d.db.QueryRow("SELECT password FROM users WHERE login=$1;", login).Scan(&u.Password); err != nil {
		return fmt.Errorf("no such user")
	}
	if password != u.Password {
		return fmt.Errorf("wrong password")
	}
	return nil
}

func (d *DB) addNewPost(imageLinks []string, text string, sessionID string) error {
	u := newUser()
	p := newPost()
	if u.ID, _, _, _, err = d.getUserInfo(sessionID); err != nil {
		return err
	}
	p.PostedOn = time.Now().Format("~~ 15:04:05 ~~ 02.01.2006 ~~")
	if err = d.db.QueryRow("INSERT INTO posts (posted_on, text, user_id) VALUES ($1, $2, $3) RETURNING id;", p.PostedOn, text, u.ID).Scan(&p.ID); err != nil {
		return err
	}
	for _, v := range imageLinks {
		if _, err = d.db.Exec("INSERT INTO images (link, post_id) VALUES ($1, $2);", v, p.ID); err != nil {
			return err
		}
	}
	return nil
}

func (d *DB) removeAPost(postID string) error {
	if _, err = d.db.Exec("DELETE FROM posts WHERE id=$1", postID); err != nil {
		return err
	}
	return nil
}

func (d *DB) removePosts(confirmationPassword string, sessionID string) error {
	u := newUser()
	if u.ID, _, u.Password, _, err = d.getUserInfo(sessionID); err != nil {
		return err
	}
	if confirmationPassword != u.Password {
		return fmt.Errorf("wrong password")
	}
	if _, err = d.db.Exec("DELETE FROM posts WHERE user_id=$1", u.ID); err != nil {
		return err
	}
	return nil
}

func (d *DB) getUserInfo(sessionID string) (userID string, login string, password string, createdOn string, err error) {
	u := newUser()
	if err = d.db.QueryRow("SELECT id, login, password, created_on FROM users WHERE session=$1", sessionID).Scan(&u.ID, &u.Login, &u.Password, &u.CreatedOn); err != nil {
		return "", "", "", "", err
	}
	return u.ID, u.Login, u.Password, u.CreatedOn, nil
}

func (d *DB) getHomePageContents(sessionID string) (login string, err error) {
	u := newUser()
	if _, u.Login, _, _, err = d.getUserInfo(sessionID); err != nil {
		return "", err
	}
	return u.Login, nil
}

func (d *DB) getBlogContents(sessionID string) (u User, err error) {
	if u.ID, u.Login, _, u.CreatedOn, err = d.getUserInfo(sessionID); err != nil {
		return u, err
	}
	favs, err := d.db.Query("SELECT favlogin FROM favourites WHERE user_id=$1", u.ID)
	if err != nil {
		return u, err
	}
	for favs.Next() {
		fp := FavPage{}
		if err = favs.Scan(&fp.Username); err != nil {
			return u, err
		}
		u.FavouritePages = append(u.FavouritePages, fp)
	}
	posts, err := d.db.Query("SELECT posts.id, posts.posted_on, posts.text FROM users INNER JOIN posts ON users.id=posts.user_id WHERE users.session=$1", sessionID)
	if err != nil {
		return u, err
	}
	u.Blog, err = d.addSQLRowsToUserStruct(&u, posts)
	if err != nil {
		return u, err
	}
	return u, nil
}

func (d *DB) getVisitContents(login string) (u User, err error) {
	if err = d.db.QueryRow("SELECT created_on FROM users where login=$1", login).Scan(&u.CreatedOn); err != nil {
		return u, err
	}
	posts, err := d.db.Query("SELECT posts.id, posts.posted_on, posts.text FROM users INNER JOIN posts ON users.id=posts.user_id WHERE users.login=$1", login)
	if err != nil {
		return u, err
	}
	blog, err := d.addSQLRowsToUserStruct(&u, posts)
	if err != nil {
		return u, err
	}
	u.Login = login
	u.Blog = blog
	return u, nil
}

func (d *DB) isPageInFavourites(visitlogin string, sessionID string) error {
	myID, _, _, _, err := d.getUserInfo(sessionID)
	if err != nil {
		return err
	}
	if err = d.db.QueryRow("SELECT user_id FROM favourites WHERE favlogin=$1 AND user_id=$2", visitlogin, myID).Scan(&myID); err != nil {
		return err
	}
	return nil
}

func (d *DB) addPageToFavourites(visitlogin string, sessionID string) error {
	myID, _, _, _, err := d.getUserInfo(sessionID)
	if err != nil {
		return err
	}
	if _, err = d.db.Exec("INSERT INTO favourites (favlogin,user_id) VALUES ($1,$2)", visitlogin, myID); err != nil {
		return err
	}
	return nil
}

func (d *DB) removePageFromFavourites(visitlogin string, sessionID string) error {
	myID, _, _, _, err := d.getUserInfo(sessionID)
	if err != nil {
		return err
	}
	if _, err = d.db.Exec("DELETE FROM favourites WHERE favlogin=$1 AND user_id=$2", visitlogin, myID); err != nil {
		return err
	}
	return nil
}

func (d *DB) addSQLRowsToUserStruct(u *User, posts *sql.Rows) (blog []Post, err error) {
	for posts.Next() {
		p := newPost()
		p.Author = u.Login
		if err = posts.Scan(&p.ID, &p.PostedOn, &p.Text); err != nil {
			return nil, err
		}
		images, err := d.db.Query("SELECT images.link FROM posts INNER JOIN images ON posts.id=images.post_id WHERE posts.id=$1", p.ID)
		if err != nil {
			return nil, err
		}
		for images.Next() {
			i := Image{}
			if err = images.Scan(&i.Link); err != nil {
				return nil, err
			}
			p.Images = append(p.Images, i)
		}
		u.Blog = append(u.Blog, p)
	}
	return u.Blog, nil
}

func (d *DB) searchSession(sessionID string) error {
	u := newUser()
	if err = d.db.QueryRow("SELECT login FROM users WHERE session=$1", sessionID).Scan(&u.Login); err != nil {
		return err
	}
	return nil
}

func (d *DB) addSession(login, sessionID string) error {
	_, err = d.db.Exec("UPDATE users SET session=$1 WHERE login=$2", sessionID, login)
	if err != nil {
		return err
	}
	return nil
}
