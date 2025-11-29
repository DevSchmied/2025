package server

import "github.com/gin-gonic/gin"

type Server struct {
	router  *gin.Engine
	address string
}

func NewServer(addr string) *Server {
	r := gin.Default()
	return &Server{
		router:  r,
		address: addr,
	}
}

func (s *Server) Start() error {

	return s.router.Run(s.address)
}
