package db

import (
	"fmt"
	"strings"
	"time"
)

type user struct {
	login, password, content string
}

func newUser() *user {
	return &user{}
}

func (s *DB) RegisterUser(login, password string) error {
	_, err := s.db.Exec("INSERT INTO users(login, password, content) VALUES ($1, $2, '');", login, password)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_pkey\"" {
			return fmt.Errorf("user already exists")
		}
		return err
	}
	return nil
}

func (s *DB) SearchUser(login, password string) (string, error) {
	u := newUser()
	row := s.db.QueryRow("SELECT * FROM users WHERE login=$1", login)
	if err := row.Scan(&u.login, &u.password, &u.content); err != nil {
		return "", fmt.Errorf("no such user")
	}
	if password != u.password {
		return "", fmt.Errorf("wrong password")
	}
	return u.content, nil
}

func (s *DB) NewBlog(login, newcnt string) (string, error) {
	updatedContent := ""
	if err := s.db.QueryRow("SELECT content from users WHERE login=$1", login).Scan(&updatedContent); err != nil {
		return updatedContent, err
	}
	errContent := updatedContent
	updatedContent += time.Now().Format(login+" posted ~~ 15:04:05 ~~ 02.01.2006 ~~") + "\n" + newcnt + "\n"
	_, err := s.db.Exec("UPDATE users SET content=$1 WHERE login=$2", updatedContent, login)
	if err != nil {
		return errContent, err
	}
	return updatedContent, nil
}

func (s *DB) RemovePosts(login string) error {
	_, err := s.db.Exec("UPDATE users SET content='' WHERE login=$1", login)
	return err
}

func (s *DB) RemoveAPost(index int, login string) (string, error) {
	allPrevCnt := ""
	allNewCnt := ""
	if err := s.db.QueryRow("SELECT content from users WHERE login=$1", login).Scan(&allPrevCnt); err != nil {
		return allPrevCnt, err
	}
	if err := s.RemovePosts(login); err != nil {
		return allPrevCnt, err
	}
	allContentSlice := strings.Split(allPrevCnt, "\n")
	for i, v := range allContentSlice {
		if i != index && index-i != 1 && v != "" {
			allNewCnt += v + "\n"
		}
	}
	_, err := s.db.Exec("UPDATE users SET content=$1 WHERE login=$2 RETURNING content", allNewCnt, login)
	if err != nil {
		return allPrevCnt, err
	}
	return allNewCnt, nil
}
