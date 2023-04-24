package main

import (
	"github.com/glebpepega/blog/server"
)

func main() {
	s := server.New()
	s.Start()
}
