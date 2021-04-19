package goftx

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient_GetServerTime(t *testing.T) {
	ftx := New()

	serverTime, err := ftx.GetServerTime()
	require.NoError(t, err)
	fmt.Println(serverTime.Sub(time.Now().UTC()))
}

func TestClient_Ping(t *testing.T) {
	t.Run("com", func(t *testing.T) {
		ftx := New(WithFTXUS())
		err := ftx.Ping()
		require.NoError(t, err)
	})

	t.Run("us", func(t *testing.T) {
		ftx := New()
		err := ftx.Ping()
		require.NoError(t, err)
	})
}
