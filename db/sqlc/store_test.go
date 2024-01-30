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
	}
}
