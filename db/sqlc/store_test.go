package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStore_TransferTx(t *testing.T) {
	store := NewStore(testDB)

	createdFromAccount := createRandomAccount(t)
	createdToAccount := createRandomAccount(t)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	n := 6
	amount := int64(10)
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: createdFromAccount.ID,
				ToAccountID:   createdToAccount.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer

		require.Equal(t, createdFromAccount.ID, transfer.FromAccountID)
		require.Equal(t, createdToAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, createdFromAccount.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, createdToAccount.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, createdFromAccount.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, createdToAccount.ID, toAccount.ID)

		// check balance
		diff1 := createdFromAccount.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - createdToAccount.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0 && diff2 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccount1, err := testQueries.GetAccount(context.Background(), createdFromAccount.ID)
	require.NoError(t, err)
	updatedAccount2, err := testQueries.GetAccount(context.Background(), createdToAccount.ID)
	require.NoError(t, err)

	require.Equal(t, createdFromAccount.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, createdToAccount.Balance+int64(n)*amount, updatedAccount2.Balance)
}
