package server

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type DB struct {
	db *sql.DB
}

func NewDB() *DB {
	return &DB{}
}

func (s *DB) Start() {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=1510 dbname=blog sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	s.db = db
}
