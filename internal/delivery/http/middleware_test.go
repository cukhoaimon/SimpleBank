package http

import (
	"github.com/cukhoaimon/SimpleBank/pkg/token"
	"github.com/cukhoaimon/SimpleBank/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_authMiddleware(t *testing.T) {
	tests := []struct {
		name          string
		setupAuth     func(*testing.T, *http.Request, token.Maker)
		checkResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "200 OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, AuthorizationTypeBearer, "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "401 - Authorization not provide",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "401 invalid authorization header format",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				request.Header.Set("hehe", "hehe")
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "401 unsupported authorization type ",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "sieu cap vo dich", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "401 invalid token",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// provide JWT token, but server is using Paseto token => invalid token
				jwtMaker, err := token.NewJWTMaker(utils.RandomString(32))
				require.Nil(t, err)
				addAuthorization(t, request, jwtMaker, AuthorizationTypeBearer, "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			config := utils.Config{
				TokenDuration:     15 * time.Minute,
				TokenSymmetricKey: utils.RandomString(32),
			}

			pasetoMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
			require.Nil(t, err)

			handler := &Handler{
				Store:      nil,
				Config:     config,
				TokenMaker: pasetoMaker,
			}

			authPath := "/auth"
			handler.Router = gin.Default()
			handler.Router.GET(
				authPath,
				AuthMiddleware(handler.TokenMaker),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.Nil(t, err)

			tc.setupAuth(t, request, handler.TokenMaker)
			handler.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
