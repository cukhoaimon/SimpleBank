package gapi

import (
	"github.com/cukhoaimon/SimpleBank/internal/delivery/grpc/pb"
	db "github.com/cukhoaimon/SimpleBank/internal/usecase/sqlc"
	"github.com/cukhoaimon/SimpleBank/pkg/token"
	"github.com/cukhoaimon/SimpleBank/utils"
)

type Handler struct {
	pb.UnimplementedSimpleBankServer
	Config     utils.Config
	TokenMaker token.Maker
	Store      db.Store
}
