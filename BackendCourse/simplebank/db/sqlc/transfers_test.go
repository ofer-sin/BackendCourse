package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func CreateRandomTransfer(t *testing.T, fromAccountID, toAccountID int64) Transfers {
	arg := CreateTransferParams{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        util.RandomMoney(),
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
	from_account := CreateRandomAccount(t)
	to_account := CreateRandomAccount(t)

	CreateRandomTransfer(t, from_account.ID, to_account.ID)
}

func TestGetTransfer(t *testing.T) {
	from_account := CreateRandomAccount(t)
	to_account := CreateRandomAccount(t)

	transfer1 := CreateRandomTransfer(t, from_account.ID, to_account.ID)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1, transfer2)
}

func TestListTransfers(t *testing.T) {
	from_account := CreateRandomAccount(t)
	to_account := CreateRandomAccount(t)

	for i := 0; i < 10; i++ {
		CreateRandomTransfer(t, from_account.ID, to_account.ID)
	}

	for pages := 0; pages < 2; pages++ {
		arg := ListTransfersParams{
			FromAccountID: from_account.ID,
			ToAccountID:   to_account.ID,
			Limit:         5,
			Offset:        5 * int32(pages),
		}
		transfers, err := testQueries.ListTransfers(context.Background(), arg)
		require.NoError(t, err)
		require.Len(t, transfers, 5)

		for _, transfer := range transfers {
			require.NotEmpty(t, transfer)
		}
	}
}
