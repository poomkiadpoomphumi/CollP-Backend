package main

import (
	"net/http"
	"collp-backend/router"
)

type Server struct {
	router http.Handler
}

func NewServer() *Server {
	r := router.ThemesRouter()
	return &Server{
		router: r,
	}
}

func (s *Server) Run(addr string) {
	http.ListenAndServe(addr, s.router)
}
