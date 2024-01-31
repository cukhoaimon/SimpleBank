package api

import (
	db "github.com/cukhoaimon/SimpleBank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer will return new HTTP server and setup router
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// routing here
	router.GET("/api/v1/account", server.listAccount)
	router.POST("/api/v1/account", server.createAccount)
	router.GET("/api/v1/account/:id", server.getAccount)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
