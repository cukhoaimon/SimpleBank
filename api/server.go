package api

import (
	"fmt"
	db "github.com/cukhoaimon/SimpleBank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer will return new HTTP server and setup router
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			fmt.Printf("err %v", err.Error())
			return nil
		}
	}

	// routing here
	router.POST("/api/v1/user", server.createUser)

	router.GET("/api/v1/account", server.listAccount)
	router.GET("/api/v1/account/:id", server.getAccount)
	router.POST("/api/v1/account", server.createAccount)

	router.POST("/api/v1/transfer", server.createTransfer)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
