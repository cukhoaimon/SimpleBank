package http

import (
	"github.com/cukhoaimon/SimpleBank/internal/delivery/http"
	db "github.com/cukhoaimon/SimpleBank/internal/usecase/sqlc"
	"github.com/cukhoaimon/SimpleBank/utils"
	"log"
)

type Server struct {
	Handler *http.Handler
}

// NewServer will return new HTTP server and setup Router
func NewServer(store db.Store, config utils.Config) (*Server, error) {
	handler, err := http.NewHandler(store, config)
	if err != nil {
		return nil, err
	}
	
	return &Server{Handler: handler}, nil
}

func (server *Server) Start(address string) error {
	return server.Handler.Router.Run(address)
}

// Run will run Gin server to serve http request
func Run(store db.Store, config utils.Config) {
	server, err := NewServer(store, config)
	if err != nil {
		log.Fatal("Cannot create server")
	}

	if err = server.Start(config.HttpServerAddress); err != nil {
		log.Fatal("Cannot start server")
	}
}
