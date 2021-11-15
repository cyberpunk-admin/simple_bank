package db

import (
	"context"
	"github.com/simplebank/db/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func CreateRandomEntry(t *testing.T, A Account) Entry{
	arg := CreateEntryParams{
		AccountID: A.ID,
		Amount: util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	return entry
}

func TestCreateEntry(t *testing.T) {
	account := CreateRandomAccount(t)
	CreateRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account1 := CreateRandomAccount(t)
	entry1 := CreateRandomEntry(t, account1)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, account1.ID, entry2.AccountID)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)

}

func TestListEntry(t *testing.T) {
	account := CreateRandomAccount(t)
	for i:=0; i < 10; i++ {
		CreateRandomEntry(t, account)
	}
	arg := ListEntryParams{
		AccountID: account.ID,
		Limit: 10,
		Offset: 0,
	}
	entries, err := testQueries.ListEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entries)

	require.Len(t, entries, 10)
	for _, et := range entries {
		require.NotEmpty(t, et)
		require.Equal(t, account.ID, et.AccountID)
	}
}

