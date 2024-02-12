package http

import (
	"fmt"
	"github.com/cukhoaimon/SimpleBank/internal/delivery/http/middleware"
	db "github.com/cukhoaimon/SimpleBank/internal/usecase/sqlc"
	"github.com/cukhoaimon/SimpleBank/pkg/token"
	"github.com/cukhoaimon/SimpleBank/utils"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"net/http"
)

func NewTestHandler(t *testing.T, store db.Store) *Handler {
	config := utils.Config{
		TokenDuration:     15 * time.Minute,
		TokenSymmetricKey: utils.RandomString(32),
	}

	handler, err := NewHandler(store, config)
	require.Nil(t, err)

	return handler
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	username string,
	duration time.Duration,
) {
	accessToken, payload, err := tokenMaker.CreateToken(username, duration)
	require.Nil(t, err)
	require.NotEmpty(t, payload)
	require.NotEmpty(t, accessToken)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, accessToken)
	request.Header.Set(middleware.AuthorizationHeaderKey, authorizationHeader)
}
