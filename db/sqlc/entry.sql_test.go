package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    10,
	}

	entry, err := testQuery.CreateEntry(context.Background(), arg)

	require.Nil(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	return entry
}

func TestQueries_CreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestQueries_GetEntry(t *testing.T) {
	account := createRandomAccount(t)
	want := createRandomEntry(t, account)

	have, err := testQuery.GetEntry(context.Background(), want.ID)

	require.Nil(t, err)
	require.NotEmpty(t, have)
	require.Equal(t, want.ID, have.ID)
	require.Equal(t, want.AccountID, have.AccountID)
	require.Equal(t, want.Amount, have.Amount)
	require.WithinDuration(t, want.CreatedAt, have.CreatedAt, time.Second)
}

func TestQueries_ListEntries(t *testing.T) {
	account := createRandomAccount(t)
	var entries []Entry
	for i := 0; i < 10; i++ {
		entries = append(entries, createRandomEntry(t, account))
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	have, err := testQuery.ListEntries(context.Background(), arg)

	require.Nil(t, err)
	require.Len(t, have, 5)

	for i, entry := range have {
		require.NotEmpty(t, entry)
		require.Equal(t, account.ID, entry.AccountID)
		require.Equal(t, entries[i+int(arg.Offset)].Amount, entry.Amount)
		require.WithinDuration(t, entries[i+int(arg.Offset)].CreatedAt, entry.CreatedAt, time.Second)
	}
}
