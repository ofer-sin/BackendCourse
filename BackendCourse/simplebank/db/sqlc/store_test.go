package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	// fmt.Println("Before Tx: ", account1, account2)

	// run concurrent transfer transactions
	n := 5
	amount := int64(10)

	// create channels to handle errors and results
	errsChannel := make(chan error)
	resultsChannel := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		// give every goroutine a name
		// this is just for debugging purposes
		// txName := fmt.Sprintf("tx-%d", i+1)
		go func() {
			// add the transaction name to the context
			// ctx := context.WithValue(context.Background(), txKey, txName)
			ctx := context.Background()
			transactionResult, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errsChannel <- err
			resultsChannel <- transactionResult
		}()
	}

	existed := make(map[int]bool) // to check that the transfer is not duplicated
	for i := 0; i < n; i++ {
		// get the error and result from the channels
		// this will block until a value is sent to the channel
		err := <-errsChannel
		require.NoError(t, err)

		transactionResult := <-resultsChannel
		require.NotEmpty(t, transactionResult)

		// fmt.Println("Transaction FromAccount: ", transactionResult.FromAccount)
		// fmt.Println("Transaction ToAccount: ", transactionResult.ToAccount)

		// check trtansfer
		transfer := transactionResult.Transfer
		require.NotEmpty(t, transactionResult)
		require.NotEmpty(t, transactionResult.Transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transactionResult.Transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := transactionResult.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := transactionResult.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := transactionResult.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := transactionResult.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check account balances
		// fmt.Println("Tx Balance: ", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance // ammont deducted from account1
		diff2 := toAccount.Balance - account2.Balance   // amount added to account2
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.Equal(t, int64(diff1%amount), int64(0))

		// there are several transaction running concurrently, each will deduct
		// the same amount from account1 and add the same amount to account2
		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.False(t, existed[k]) // check that the transfer is not duplicated
		existed[k] = true
	}

	// check the final updated balances
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)

	// fmt.Println("After Tx: ", updatedAccount1.Balance, updatedAccount2.Balance)
}

func TestTransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	// fmt.Println("Before Tx: ", account1, account2)

	// run concurrent transfer transactions
	n := 10
	amount := int64(10)

	// create channels to handle errors and results
	errsChannel := make(chan error)

	for i := 0; i < n; i++ {
		// give every goroutine a name
		// this is just for debugging purposes
		// txName := fmt.Sprintf("tx-%d", i+1)
		// switch between the two accounts in each iteration
		// this will create a deadlock if the two transactions are running concurrently
		// and trying to lock the same two accounts in different order
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			// add the transaction name to the context
			// ctx := context.WithValue(context.Background(), txKey, txName)
			ctx := context.Background()
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errsChannel <- err
		}()
	}

	for i := 0; i < n; i++ {
		// get the error and result from the channels
		// this will block until a value is sent to the channel
		err := <-errsChannel
		require.NoError(t, err)

	}

	// check the final updated balances
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)

	// fmt.Println("After Tx: ", updatedAccount1.Balance, updatedAccount2.Balance)
}
