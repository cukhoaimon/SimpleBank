package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/cukhoaimon/SimpleBank/internal/delivery/http/mock"
	db "github.com/cukhoaimon/SimpleBank/internal/usecase/sqlc"

	"github.com/cukhoaimon/SimpleBank/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func EqCreateUserPasswordParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqUserPasswordMatcher{arg, password}
}

type eqUserPasswordMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqUserPasswordMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := utils.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqUserPasswordMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func TestServer_createUser(t *testing.T) {
	user, password := randomUser(t)

	wantResponse := userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "201 Created",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"email":     user.Email,
				"full_name": user.FullName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:       user.Username,
					HashedPassword: user.HashedPassword,
					FullName:       user.FullName,
					Email:          user.Email,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserPasswordParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireEqualUserBody(t, recorder.Body, wantResponse)
			},
		},
		{
			name: "401 Bad request - Body mismatch",
			body: gin.H{
				"ten_nguoi_dung": user.Username,
				"mat_khau":       password,
				"email":          user.Email,
				"ten_day_du":     user.FullName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "403 Forbidden - unique_violation",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"email":     user.Email,
				"full_name": user.FullName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: pq.ErrorCode("23505")})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			// build stubs
			tc.buildStubs(store)

			// start test handler and send request
			handler := NewTestHandler(t, store)
			recorder := httptest.NewRecorder()

			url := "/api/v1/user"

			data, err := json.Marshal(tc.body)
			require.Nil(t, err)

			// send request
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.Nil(t, err)

			handler.Router.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func TestServer_loginUser(t *testing.T) {
	user, password := randomUser(t)
	userResponse := makeUserResponse(user)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "201 Created",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1)

				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var haveLoginResponse loginResponse
				data, err := io.ReadAll(recorder.Body)
				require.Nil(t, err)

				err = json.Unmarshal(data, &haveLoginResponse)
				require.Nil(t, err)

				require.NotEmpty(t, haveLoginResponse.AccessToken)
				require.Equal(t, userResponse, haveLoginResponse.User)
			},
		},
		{
			name: "400 Bad request - error when binding json",
			body: gin.H{
				"ten_nguoi_dung": user.Username,
				"mat_khau":       password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "404 Not found - sql no row error",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "500 Internal handler error",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "401 Unauthorized - wrong password",
			body: gin.H{
				"username": user.Username,
				"password": utils.RandomString(10),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			// build stubs
			tc.buildStubs(store)

			// start test handler and send request
			handler := NewTestHandler(t, store)
			recorder := httptest.NewRecorder()

			url := "/api/v1/user/login"

			data, err := json.Marshal(tc.body)
			require.Nil(t, err)

			// send request
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.Nil(t, err)

			handler.Router.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = utils.RandomString(10)
	hashedPassword, err := utils.HashPassword(password)
	require.Nil(t, err)

	user = db.User{
		FullName:       utils.RandomOwner(),
		Username:       utils.RandomOwner(),
		HashedPassword: hashedPassword,
		Email:          utils.RandomEmail(),
	}

	return
}

func requireEqualUserBody(t *testing.T, body *bytes.Buffer, want userResponse) {
	var have db.User

	data, err := io.ReadAll(body)
	require.Nil(t, err)

	err = json.Unmarshal(data, &have)
	require.Nil(t, err)

	haveResponse := userResponse{
		Username:          have.Username,
		FullName:          have.FullName,
		Email:             have.Email,
		PasswordChangedAt: have.PasswordChangedAt,
		CreatedAt:         have.CreatedAt,
	}

	require.Equal(t, want, haveResponse)
}
