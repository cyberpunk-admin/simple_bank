package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTX(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	fmt.Println(">> begin: ", account1.Balance, account2.Balance)
	// run n concurrent transfer transcations
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTXResult)

	for i:= 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTXParams{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount: amount,
			})
			errs <- err
			results <- result
		}()
	}
	// check result
	existed := make(map[int]bool)

	for i:= 0; i < n; i++ {
		err := <- errs
		require.NoError(t, err)
		result := <- results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entres
		FromEntry := result.FromEntry
		require.NotEmpty(t, FromEntry)
		require.Equal(t, account1.ID, FromEntry.AccountID)
		require.Equal(t, -amount, FromEntry.Amount)
		require.NotZero(t, FromEntry.ID)
		require.NotZero(t, FromEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), FromEntry.ID)
		require.NoError(t, err)

		ToEntry := result.ToEntry
		require.NotEmpty(t, ToEntry)
		require.Equal(t, account2.ID, ToEntry.AccountID)
		require.Equal(t, amount, ToEntry.Amount)
		require.NotZero(t, ToEntry.ID)
		require.NotZero(t, ToEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), ToEntry.ID)
		require.NoError(t, err)

		// check account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check balance
		fmt.Println(">> tx: ", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1 % amount == 0)

		k := int(diff1 / amount)
		require.True(t, k > 0 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check final update balance
	updateFromAccount, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateFromAccount)
	require.True(t, updateFromAccount.Balance + int64(n) * amount == account1.Balance)

	updateToAccount, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateToAccount)
	require.True(t, updateToAccount.Balance - int64(n) * amount == account2.Balance)

	fmt.Println(">> after:", updateFromAccount.Balance, updateToAccount.Balance)
}

func TestTransferTXDeadLock(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	fmt.Println(">> begin: ", account1.Balance, account2.Balance)
	// run n concurrent transfer transcations
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i:= 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID
		if i % 2 == 1 {
			fromAccountID, toAccountID = toAccountID, fromAccountID
		}

		go func() {
			//ctx := context.WithValue(context.Background(), txKey, txName)
			_, err := store.TransferTx(context.Background(), TransferTXParams{
				FromAccountID: fromAccountID,
				ToAccountID: toAccountID,
				Amount: amount,
			})
			errs <- err
		}()
	}
	// check result
	for i:= 0; i < n; i++{
		err := <- errs
		require.NoError(t, err)
	}
	// check final update balance
	updateFromAccount, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateFromAccount)
	require.Equal(t, account1.Balance,updateFromAccount.Balance)


	updateToAccount, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateToAccount)

	require.Equal(t, account2.Balance, updateToAccount.Balance)

	fmt.Println(">> after:", updateFromAccount.Balance, updateToAccount.Balance)
}