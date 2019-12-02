package main

import (
	"net/http"
)

func (s *server) routes() {
	s.router.ServeFiles("/*filepath", http.Dir("./ui/"))
	s.POST("/irccip/key", s.handleKeyPress())
}

func (s *server) POST(path string, h http.HandlerFunc) {
	s.router.HandlerFunc(http.MethodPost, path, h)
}
