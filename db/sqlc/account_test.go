package db

import (
	"context"
	"testing"

	"github.com/GK16/miniBank/util"
	"github.com/stretchr/testify/require"
)

func TestCreatAccount(t *testing.T){
	arg := CreateAccountParams{
		Owner: util.RandOwner(),
		Balance: util.RandMoney(),
		Currency: util.RandCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}