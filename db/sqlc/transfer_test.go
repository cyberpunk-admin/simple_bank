package db

import (
	"context"
	"github.com/simplebank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func CreateRandomTransfer(t *testing.T, from Account, to Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: from.ID,
		ToAccountID:   to.ID,
		Amount:        util.RandomAmount(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	return transfer
}

func TestCreateTransfer(t *testing.T) {
	accountFrom := CreateRandomAccount(t)
	accountTO := CreateRandomAccount(t)
	CreateRandomTransfer(t, accountFrom, accountTO)
}

func TestGetTransfer(t *testing.T) {
	accountFrom := CreateRandomAccount(t)
	accountTO := CreateRandomAccount(t)
	transfer1 := CreateRandomTransfer(t, accountFrom, accountTO)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)

}

func TestListTransfer(t *testing.T) {
	accountFrom := CreateRandomAccount(t)
	accountTO := CreateRandomAccount(t)

	for i := 0; i < 5; i++ {
		CreateRandomTransfer(t, accountFrom, accountTO)
		CreateRandomTransfer(t, accountTO, accountFrom)
	}

	arg := ListTransferParams{
		FromAccountID: accountFrom.ID,
		ToAccountID:   accountFrom.ID,
		Limit:         5,
		Offset:        5,
	}
	transfers, err := testQueries.ListTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	require.Len(t, transfers, 5)

	for _, tran := range transfers {
		require.NotEmpty(t, tran)
		require.True(t, tran.FromAccountID == arg.FromAccountID || tran.ToAccountID == arg.FromAccountID)
	}
}
