package goftx

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grishinsana/goftx/models"
)

func TestFills_GetFills(t *testing.T) {
	_ = godotenv.Load()

	ftx := New(
		WithAuth(os.Getenv("FTX_KEY"), os.Getenv("FTX_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	require.NoError(t, err)

	market := "ETH/BTC"

	fills, err := ftx.Fills.GetFills(&models.GetFillsParams{
		Market: &market,
	})
	assert.NoError(t, err)
	assert.NotNil(t, fills)
}
