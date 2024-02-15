package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/cukhoaimon/SimpleBank/internal/delivery/http/mock"
	db "github.com/cukhoaimon/SimpleBank/internal/usecase/sqlc"
	"github.com/cukhoaimon/SimpleBank/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestServer_getAccount(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccountWithUser(user)

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "200 OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)). /* call getAccount in any context */
					Times(1).                                        /* how many times call function */
					Return(account, nil)                             /* must fit with function */
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "Not found",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "500 Internal server error",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "400 Bad request",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			// start test server and send request
			handler := NewTestHandler(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/account/%d", tc.accountID)

			// send request
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.Nil(t, err)

			addAuthorization(t, request, handler.TokenMaker, AuthorizationTypeBearer, user.Username, time.Minute)
			handler.Router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func TestServer_createAccount(t *testing.T) {
	arg := db.CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  0,
		Currency: utils.RandomCurrency(),
	}

	wantAccount := db.Account{
		ID:       utils.RandomInt(1, 1000),
		Owner:    arg.Owner,
		Balance:  arg.Balance,
		Currency: arg.Currency,
	}

	tests := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "201 Created",
			body: gin.H{
				"owner":    arg.Owner,
				"balance":  arg.Balance,
				"currency": arg.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(wantAccount, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, wantAccount)
			},
		},
		{
			name: "400 Bad request",
			body: gin.H{
				"chu_so_huu": arg.Owner,
				"so_du":      arg.Balance,
				"tien_te":    arg.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "500 Internal error",
			body: gin.H{
				"owner":    arg.Owner,
				"balance":  arg.Balance,
				"currency": arg.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			// build stubs
			tc.buildStubs(store)

			// start test server and send request
			server := NewTestHandler(t, store)
			recorder := httptest.NewRecorder()

			url := "/api/v1/account"

			data, err := json.Marshal(tc.body)
			require.Nil(t, err)

			// send request
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
			require.Nil(t, err)

			addAuthorization(t, request, server.TokenMaker, AuthorizationTypeBearer, arg.Owner, time.Minute)

			server.Router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func TestServer_listAccount(t *testing.T) {
	user, _ := randomUser(t)
	n := 10
	accounts := make([]db.Account, n)

	for i := 0; i < n; i++ {
		accounts = append(accounts, randomAccountWithUser(user))
	}

	arg := listAccountRequest{
		Owner:    user.Username,
		PageID:   1,
		PageSize: 5,
	}

	tests := []struct {
		name          string
		query         listAccountRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "200 OK",
			query: arg,
			buildStubs: func(store *mockdb.MockStore) {
				offset := (arg.PageID - 1) * arg.PageSize

				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(db.ListAccountsParams{
						Owner:  user.Username,
						Limit:  arg.PageSize,
						Offset: offset,
					})).
					Times(1).
					Return(accounts[offset:offset+arg.PageSize], nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				// check body
				data, err := io.ReadAll(recorder.Body)
				require.Nil(t, err)

				var gotAccounts []db.Account
				err = json.Unmarshal(data, &gotAccounts)
				require.Nil(t, err)

				offset := int((arg.PageID - 1) * arg.PageSize)

				for i, gotAccount := range gotAccounts {
					require.Equal(t, accounts[offset+i], gotAccount)
				}
			},
		},
		{
			name: "400 Bad request",
			query: listAccountRequest{
				PageID:   0,
				PageSize: 0,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "500 Internal server",
			query: arg,
			buildStubs: func(store *mockdb.MockStore) {
				offset := (arg.PageID - 1) * arg.PageSize

				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(db.ListAccountsParams{
						Owner:  user.Username,
						Offset: offset,
						Limit:  arg.PageSize,
					})).
					Times(1).
					Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			// build stubs
			tc.buildStubs(store)

			// start test server and send request
			server := NewTestHandler(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/account?page_id=%d&page_size=%d", tc.query.PageID, tc.query.PageSize)

			// send request
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.Nil(t, err)
			addAuthorization(t, request, server.TokenMaker, AuthorizationTypeBearer, user.Username, time.Minute)

			server.Router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       utils.RandomInt(1, 1000),
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
}

func randomAccountWithUser(user db.User) db.Account {
	return db.Account{
		ID:       utils.RandomInt(1, 1000),
		Owner:    user.Username,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.Nil(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)

	require.Nil(t, err)

	require.Equal(t, account, gotAccount)
}
