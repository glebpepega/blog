package database

import (
	"fmt"
	"strings"
	"time"
)

type user struct {
	login     string
	password  string
	content   string
	sessionID string
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

func (s *DB) SearchUser(login, password string) (content string, err error) {
	u := user{}
	if err := s.db.QueryRow("SELECT login, password, content FROM users WHERE login=$1", login).Scan(&u.login, &u.password, &u.content); err != nil {
		return "", fmt.Errorf("no such user")
	}
	if password != u.password {
		return "", fmt.Errorf("wrong password")
	}
	return u.content, nil
}

func (s *DB) NewBlog(login, newcnt string) (updatedContent string, err error) {
	if err := s.db.QueryRow("SELECT content from users WHERE login=$1", login).Scan(&updatedContent); err != nil {
		return updatedContent, err
	}
	errContent := updatedContent
	updatedContent += time.Now().Format(login+" posted ~~ 15:04:05 ~~ 02.01.2006 ~~") + "\n" + newcnt + "\n"
	_, err = s.db.Exec("UPDATE users SET content=$1 WHERE login=$2", updatedContent, login)
	if err != nil {
		return errContent, err
	}
	return updatedContent, nil
}

func (s *DB) RemoveAPost(login string, index int) (allNewCnt string, err error) {
	allPrevCnt := ""
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
	_, err = s.db.Exec("UPDATE users SET content=$1 WHERE login=$2 RETURNING content", allNewCnt, login)
	if err != nil {
		return allPrevCnt, err
	}
	return allNewCnt, nil
}

func (s *DB) RemovePosts(login string) error {
	_, err := s.db.Exec("UPDATE users SET content='' WHERE login=$1", login)
	return err
}

func (s *DB) SearchSession(sessionID string) (login string, content string, err error) {
	u := user{}
	if err := s.db.QueryRow("SELECT login, content, sessionid from users WHERE sessionid=$1", sessionID).Scan(&u.login, &u.content, &u.sessionID); err != nil {
		return "", "", err
	}
	return u.login, u.content, nil
}

func (s *DB) AddSession(login, sessionID string) error {
	_, err := s.db.Exec("UPDATE users SET sessionid=$1 WHERE login=$2", sessionID, login)
	if err != nil {
		return err
	}
	return nil
}
