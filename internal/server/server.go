package server

import (
	"2025/internal/service"
	"2025/internal/storage"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router  *gin.Engine
	address string
	storage *storage.Storage
	tasks   chan service.Task
}

func NewServer(addr string, strg *storage.Storage, tsks chan service.Task) *Server {
	r := gin.Default()
	return &Server{
		router:  r,
		address: addr,
		storage: strg,
		tasks:   tsks,
	}
}

func (s *Server) Start() error {

	return s.router.Run(s.address)
}
