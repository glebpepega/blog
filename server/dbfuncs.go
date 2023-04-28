package server

import (
	"fmt"
	"time"
)

type user struct {
	id       string
	login    string
	password string
	blog     []*Post
}

func newUser() *user {
	return &user{}
}

func newPost() *Post {
	return &Post{}
}

func (d *DB) registerUser(login, password string) error {
	if _, err = d.db.Exec("INSERT INTO users (login, password) VALUES ($1, $2);", login, password); err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_pkey\"" {
			return fmt.Errorf("user already exists")
		}
		return err
	}
	return nil
}

func (d *DB) searchUser(login, password string) error {
	u := newUser()
	if err = d.db.QueryRow("SELECT password FROM users WHERE login=$1;", login).Scan(&u.password); err != nil {
		return fmt.Errorf("no such user")
	}
	if password != u.password {
		return fmt.Errorf("wrong password")
	}
	return nil
}

func (d *DB) addNewPost(imageLinks []string, text string, sessionID string) (err error) {
	p := newPost()
	u := newUser()
	if u.id, _, _, err = d.getUserInfo(sessionID); err != nil {
		return err
	}
	p.Timestamp = time.Now().Format("~~ 15:04:05 ~~ 02.01.2006 ~~")
	if err = d.db.QueryRow("INSERT INTO posts (timestamp, text, user_id) VALUES ($1, $2, $3) RETURNING id;", p.Timestamp, text, u.id).Scan(&p.ID); err != nil {
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

func (d *DB) removePosts(confirmationPassword string, sessionID string) (err error) {
	u := newUser()
	if u.id, _, u.password, err = d.getUserInfo(sessionID); err != nil {
		return err
	}
	if confirmationPassword != u.password {
		return fmt.Errorf("wrong password")
	}
	if _, err = d.db.Exec("DELETE FROM posts WHERE user_id=$1", u.id); err != nil {
		return err
	}
	return nil
}

func (d *DB) getUserInfo(sessionID string) (userID string, login string, password string, err error) {
	u := newUser()
	if err = d.db.QueryRow("SELECT id, login, password FROM users WHERE session=$1", sessionID).Scan(&u.id, &u.login, &u.password); err != nil {
		return "", "", "", err
	}
	return u.id, u.login, u.password, nil
}

func (d *DB) getHomePageContents(sessionID string) (login string, err error) {
	u := newUser()
	if _, u.login, _, err = d.getUserInfo(sessionID); err != nil {
		return "", err
	}
	return u.login, nil
}

func (d *DB) getBlogContents(sessionID string) (login string, blogCnts []*Post, err error) {
	u := newUser()
	if _, u.login, _, err = d.getUserInfo(sessionID); err != nil {
		return "", nil, err
	}
	posts, err := d.db.Query("SELECT posts.id, posts.timestamp, posts.text FROM users INNER JOIN posts ON users.id=posts.user_id WHERE users.session=$1", sessionID)
	if err != nil {
		return u.login, nil, err
	}
	for posts.Next() {
		p := newPost()
		p.Author = u.login
		if err = posts.Scan(&p.ID, &p.Timestamp, &p.Text); err != nil {
			return u.login, nil, err
		}
		images, err := d.db.Query("SELECT images.link FROM posts INNER JOIN images ON posts.id=images.post_id WHERE posts.id=$1", p.ID)
		if err != nil {
			return u.login, nil, err
		}
		for images.Next() {
			i := Image{}
			if err = images.Scan(&i.Link); err != nil {
				return u.login, nil, err
			}
			p.Images = append(p.Images, i)
		}
		u.blog = append(u.blog, p)
	}
	return u.login, u.blog, nil
}

func (d *DB) searchSession(sessionID string) error {
	u := newUser()
	if err = d.db.QueryRow("SELECT login FROM users WHERE session=$1", sessionID).Scan(&u.login); err != nil {
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
