package goftx

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/grishinsana/goftx/models"
)

func TestPrepareQueryParams(t *testing.T) {
	tests := []struct {
		params   interface{}
		expected map[string]string
		err      error
	}{
		{
			params: &models.GetTradesParams{
				Limit: nil,
			},
			expected: map[string]string{},
			err:      nil,
		},
		{
			params: &models.GetTradesParams{
				Limit:     PtrInt(10),
				StartTime: PtrInt(20),
				EndTime:   PtrInt(30),
			},
			expected: map[string]string{
				"limit":      "10",
				"start_time": "20",
				"end_time":   "30",
			},
			err: nil,
		},
		{
			params: &models.GetTradesParams{
				Limit:     PtrInt(10),
				StartTime: PtrInt(20),
				EndTime:   PtrInt(0),
			},
			expected: map[string]string{
				"limit":      "10",
				"start_time": "20",
				"end_time":   "0",
			},
			err: nil,
		},
		{
			params: &models.GetHistoricalPricesParams{
				Limit: PtrInt(10),
			},
			expected: map[string]string{},
			err:      errors.New("required field: resolution"),
		},
		{
			params: &models.GetHistoricalPricesParams{
				Resolution: models.Minute,
				Limit:      PtrInt(10),
				StartTime:  PtrInt(20),
				EndTime:    PtrInt(0),
			},
			expected: map[string]string{
				"resolution": "60",
				"limit":      "10",
				"start_time": "20",
				"end_time":   "0",
			},
			err: nil,
		},
	}

	for i, test := range tests {
		msg := fmt.Sprintf("test #%d", i+1)
		result, err := PrepareQueryParams(test.params)
		if err != nil {
			require.NotNil(t, test.err)
			require.Equal(t, test.err.Error(), err.Error(), msg)
		}
		require.Len(t, result, len(test.expected), msg)
		for k, v := range test.expected {
			value, ok := result[k]
			require.Equal(t, true, ok, msg)
			require.Equal(t, v, value, msg)
		}
	}
}

func PtrInt(i int) *int {
	return &i
}
