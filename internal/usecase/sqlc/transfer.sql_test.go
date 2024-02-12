package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/cukhoaimon/SimpleBank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, fromAccount, toAccount Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        utils.RandomMoney(),
	}

	transfer, err := testQuery.CreateTransfer(context.Background(), arg)

	require.Nil(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	return transfer
}

func TestQueries_CreateTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	createRandomTransfer(t, fromAccount, toAccount)
}

func TestQueries_GetTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	want := createRandomTransfer(t, fromAccount, toAccount)

	have, err := testQuery.GetTransfer(context.Background(), want.ID)

	require.Nil(t, err)
	require.NotEmpty(t, have)
	require.Equal(t, want.ID, have.ID)
	require.Equal(t, want.FromAccountID, have.FromAccountID)
	require.Equal(t, want.ToAccountID, have.ToAccountID)
	require.Equal(t, want.Amount, have.Amount)
	require.WithinDuration(t, want.CreatedAt, have.CreatedAt, time.Second)
}

func TestQueries_ListTransfers(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	var transfers []Transfer

	for i := 0; i < 10; i++ {
		transfer := createRandomTransfer(t, fromAccount, toAccount)
		transfers = append(transfers, transfer)
	}

	arg := ListTransfersParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Limit:         5,
		Offset:        5,
	}

	have, err := testQuery.ListTransfers(context.Background(), arg)

	require.Nil(t, err)
	require.Len(t, have, int(arg.Limit))

	for i, transfer := range have {
		require.NotEmpty(t, transfer)
		require.Equal(t, transfers[i+int(arg.Offset)].ID, transfer.ID)
		require.Equal(t, transfers[i+int(arg.Offset)].FromAccountID, transfer.FromAccountID)
		require.Equal(t, transfers[i+int(arg.Offset)].ToAccountID, transfer.ToAccountID)
		require.Equal(t, transfers[i+int(arg.Offset)].Amount, transfer.Amount)
		require.WithinDuration(t, transfers[i+int(arg.Offset)].CreatedAt, transfer.CreatedAt, time.Second)
	}
}
