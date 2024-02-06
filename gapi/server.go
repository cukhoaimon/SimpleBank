package gapi

import (
	db "github.com/cukhoaimon/SimpleBank/db/sqlc"
	"github.com/cukhoaimon/SimpleBank/pb"
	"github.com/cukhoaimon/SimpleBank/token"
	"github.com/cukhoaimon/SimpleBank/utils"
)

// Server serves gRPC request
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     utils.Config
	tokenMaker token.Maker
	store      db.Store
}

// NewServer will return new gRPC server
func NewServer(store db.Store, config utils.Config) (*Server, error) {
	maker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		store:      store,
		tokenMaker: maker,
		config:     config,
	}

	return server, nil
}
