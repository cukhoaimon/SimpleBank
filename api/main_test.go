package api

import (
	db "github.com/cukhoaimon/SimpleBank/db/sqlc"
	"github.com/cukhoaimon/SimpleBank/token"
	"github.com/cukhoaimon/SimpleBank/utils"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := utils.Config{
		TokenDuration:     15 * time.Minute,
		TokenSymmetricKey: utils.RandomString(32),
	}

	pasetoMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	require.Nil(t, err)

	server, err := NewServer(store, config)
	require.Nil(t, err)

	server.tokenMaker = pasetoMaker
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
