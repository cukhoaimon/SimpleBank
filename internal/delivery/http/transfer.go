package http

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/cukhoaimon/SimpleBank/internal/delivery/http/middleware"
	db "github.com/cukhoaimon/SimpleBank/internal/usecase/sqlc"
	"github.com/cukhoaimon/SimpleBank/pkg/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,min=1"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (handler *Handler) createTransfer(ctx *gin.Context) {
	var req transferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !handler.validOwner(ctx, req.FromAccountID) {
		return
	}

	if !handler.validAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}

	if !handler.validAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	account, err := handler.Store.TransferTxAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

func (handler *Handler) validOwner(ctx *gin.Context, accountID int64) bool {
	account, err := handler.Store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err = errors.New("from_account is not belong to the authorized user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return false
	}

	return true
}

func (handler *Handler) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := handler.Store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] mismatch: %s vs %s ", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}
