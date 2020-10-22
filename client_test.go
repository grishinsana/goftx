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
