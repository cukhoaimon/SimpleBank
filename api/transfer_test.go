package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	mockdb "github.com/cukhoaimon/SimpleBank/db/mock"
	db "github.com/cukhoaimon/SimpleBank/db/sqlc"
	"github.com/cukhoaimon/SimpleBank/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServer_createTransfer(t *testing.T) {
	fromAccount := db.Account{
		ID:      utils.RandomInt(0, 100),
		Owner:   utils.RandomOwner(),
		Balance: utils.RandomInt(100, 500),

		Currency: utils.VND,
	}

	toAccount := db.Account{
		ID:       utils.RandomInt(0, 100),
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.VND,
	}

	thirdAccount := db.Account{
		ID:       utils.RandomInt(0, 100),
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.EUR,
	}

	currency := utils.VND
	amount := int64(10)

	transferResult := db.TransferTxResult{
		Transfer: db.Transfer{
			ID:            utils.RandomInt(0, 100),
			FromAccountID: fromAccount.ID,
			ToAccountID:   toAccount.ID,
			Amount:        amount,
		},
		FromAccount: fromAccount,
		ToAccount:   toAccount,
		FromEntry: db.Entry{
			ID:        utils.RandomInt(0, 100),
			AccountID: fromAccount.ID,
			Amount:    -amount,
		},
		ToEntry: db.Entry{
			ID:        utils.RandomInt(0, 100),
			AccountID: toAccount.ID,
			Amount:    amount,
		},
	}

	tests := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "201 created",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).
					Times(2).
					Return(fromAccount, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).
					Times(1).
					Return(toAccount, nil)

				store.EXPECT().
					TransferTxAccount(gomock.Any(), gomock.Eq(db.TransferTxParams{
						FromAccountID: fromAccount.ID,
						ToAccountID:   toAccount.ID,
						Amount:        amount,
					})).Times(1).
					Return(transferResult, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				var haveResultTx db.TransferTxResult

				data, err := io.ReadAll(recorder.Body)
				require.Nil(t, err)

				err = json.Unmarshal(data, &haveResultTx)

				require.Equal(t, transferResult, haveResultTx)
			},
		},
		{
			name: "400 Bad request - binding json error",
			body: gin.H{
				"tu_tai_khoan_id":  fromAccount.ID,
				"den_tai_khoan_id": toAccount.ID,
				"so_luong":         amount,
				"don_vi":           currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					TransferTxAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "404 Not found - Missing record from the result set",
			body: gin.H{
				"from_account_id": 5000,
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)

				store.EXPECT().
					TransferTxAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "400 Bad request - Mismatch from the first account currency",
			body: gin.H{
				"from_account_id": fromAccount.ID, // currency is VND
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        utils.EUR,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).
					Times(2).
					Return(fromAccount, nil)

				store.EXPECT().
					TransferTxAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "400 Bad request - Mismatch from the second account currency",
			body: gin.H{
				"from_account_id": fromAccount.ID,  // currency is VND
				"to_account_id":   thirdAccount.ID, // currency is EUR
				"amount":          amount,
				"currency":        utils.VND,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).
					Times(2).
					Return(fromAccount, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(thirdAccount.ID)).
					Times(1).
					Return(thirdAccount, nil)

				store.EXPECT().
					TransferTxAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "500 Internal server error from GetAccount",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).
					Times(0)

				store.EXPECT().
					TransferTxAccount(gomock.Any(), gomock.Eq(db.TransferTxParams{
						FromAccountID: fromAccount.ID,
						ToAccountID:   toAccount.ID,
						Amount:        amount,
					})).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "500 Internal server error from TransferTxAccount",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).
					Times(2).
					Return(fromAccount, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).
					Times(1).
					Return(toAccount, nil)

				store.EXPECT().
					TransferTxAccount(gomock.Any(), gomock.Eq(db.TransferTxParams{
						FromAccountID: fromAccount.ID,
						ToAccountID:   toAccount.ID,
						Amount:        amount,
					})).Times(1).
					Return(db.TransferTxResult{}, sql.ErrConnDone)
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
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/api/v1/transfer"

			data, err := json.Marshal(tc.body)
			require.Nil(t, err)

			// send request
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
			require.Nil(t, err)

			addAuthorization(t, request, server.tokenMaker, authorizationTypeBearer, fromAccount.Owner, time.Minute)
			server.router.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder)
		})
	}
}
