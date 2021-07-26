package goftx

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestWallet_GetBalances(t *testing.T) {
	_ = godotenv.Load()

	ftx := New(
		WithAuth(os.Getenv("FTX_KEY"), os.Getenv("FTX_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	require.NoError(t, err)

	balances, err := ftx.GetBalances()
	require.Nil(t, err)
	t.Logf("%+v", balances)
	require.NotNil(t, balances)
}
