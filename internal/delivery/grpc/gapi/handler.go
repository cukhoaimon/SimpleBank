package gapi

import (
	db "github.com/cukhoaimon/SimpleBank/internal/usecase/sqlc"
	"github.com/cukhoaimon/SimpleBank/pkg/token"
	"github.com/cukhoaimon/SimpleBank/utils"
)

type Handler struct {
	Config     utils.Config
	TokenMaker token.Maker
	Store      db.Store
}
