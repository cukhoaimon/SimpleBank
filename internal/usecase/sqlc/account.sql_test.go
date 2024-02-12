package usecase

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/cukhoaimon/SimpleBank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQuery.CreateAccount(context.Background(), arg)

	require.Nil(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestQueries_CreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestQueries_GetAccount(t *testing.T) {
	want := createRandomAccount(t)
	have, err := testQuery.GetAccount(context.Background(), want.ID)

	require.Nil(t, err)

	require.NotEmpty(t, have)

	require.Equal(t, want.ID, have.ID)
	require.Equal(t, want.Owner, have.Owner)
	require.Equal(t, want.Balance, have.Balance)
	require.Equal(t, want.Currency, have.Currency)

	require.WithinDuration(t, want.CreatedAt, have.CreatedAt, time.Second)
}

func TestQueries_UpdateAccount(t *testing.T) {
	want := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      want.ID,
		Balance: utils.RandomMoney(),
	}

	have, err := testQuery.UpdateAccount(context.Background(), arg)

	require.Nil(t, err)
	require.NotEmpty(t, have)

	require.Equal(t, arg.ID, have.ID)
	require.Equal(t, arg.Balance, have.Balance)
	require.Equal(t, want.Owner, have.Owner)
	require.Equal(t, want.Currency, have.Currency)
	require.WithinDuration(t, want.CreatedAt, have.CreatedAt, time.Second)
}

func TestQueries_DeleteAccount(t *testing.T) {
	want := createRandomAccount(t)
	err := testQuery.DeleteAccount(context.Background(), want.ID)

	require.Nil(t, err)

	have, err := testQuery.GetAccount(context.Background(), want.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, have)
}

// TODO: Fix List Account in the ListAccountsParams
func TestQueries_ListAccount(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	have, err := testQuery.ListAccounts(context.Background(), arg)

	require.Nil(t, err)
	require.Len(t, have, 1)

	for _, account := range have {
		require.NotEmpty(t, account)
	}
}
