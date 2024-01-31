package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStore_TransferTxAccount(t *testing.T) {
	store := NewStore(testDB)

	wantFromAccount := createRandomAccount(t)
	wantToAccount := createRandomAccount(t)

	n := 10
	amount := int64(100)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTxAccount(context.Background(), TransferTxParams{
				FromAccountID: wantFromAccount.ID,
				ToAccountId:   wantToAccount.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)

	// validate
	for i := 0; i < n; i++ {
		err := <-errs
		require.Nil(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// Check Transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		require.Equal(t, wantFromAccount.ID, transfer.FromAccountID)
		require.Equal(t, wantToAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)

		// make sure transfer exists in db
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.Nil(t, err)

		// Check FromEntry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)

		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		require.Equal(t, wantFromAccount.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.Nil(t, err)

		// Check ToEntry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)

		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		require.Equal(t, wantToAccount.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.Nil(t, err)

		// Check account

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, wantFromAccount.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, wantToAccount.ID, toAccount.ID)

		diff1 := wantFromAccount.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - wantToAccount.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)

		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check final update
	updatedAccount1, err := testQuery.GetAccount(context.Background(), wantFromAccount.ID)
	require.Nil(t, err)
	require.Equal(t, wantFromAccount.Balance-int64(n)*amount, updatedAccount1.Balance)

	updatedAccount2, err := testQuery.GetAccount(context.Background(), wantToAccount.ID)
	require.Nil(t, err)
	require.Equal(t, wantToAccount.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func TestStore_TransferTxAccountDeadLock(t *testing.T) {
	store := NewStore(testDB)

	wantFromAccount := createRandomAccount(t)
	wantToAccount := createRandomAccount(t)

	n := 10
	amount := int64(100)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := wantFromAccount.ID
		toAccountID := wantToAccount.ID

		if i%2 == 1 {
			fromAccountID = wantToAccount.ID
			toAccountID = wantFromAccount.ID
		}

		go func() {
			_, err := store.TransferTxAccount(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountId:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	// validate
	for i := 0; i < n; i++ {
		err := <-errs
		require.Nil(t, err)
	}

	// check final update
	updatedAccount1, err := testQuery.GetAccount(context.Background(), wantFromAccount.ID)
	require.Nil(t, err)

	updatedAccount2, err := testQuery.GetAccount(context.Background(), wantToAccount.ID)
	require.Nil(t, err)

	require.Equal(t, wantFromAccount.Balance, updatedAccount1.Balance)
	require.Equal(t, wantToAccount.Balance, updatedAccount2.Balance)
}
