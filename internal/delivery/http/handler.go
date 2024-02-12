package http

import (
	db "github.com/cukhoaimon/SimpleBank/internal/usecase/sqlc"
	"github.com/cukhoaimon/SimpleBank/pkg/token"
	"github.com/cukhoaimon/SimpleBank/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	Store      db.Store
	Config     utils.Config
	TokenMaker token.Maker
	Router     *gin.Engine
}

func NewHandler(store db.Store, config utils.Config) (*Handler, error) {
	maker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	handler := &Handler{
		Store:      store,
		TokenMaker: maker,
		Config:     config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			return nil, err
		}
	}

	handler.SetupRouter()
	return handler, nil
}
