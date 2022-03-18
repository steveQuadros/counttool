package scanner

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

type testCounter struct {
	data map[string]int
}

func newTestCounter() *testCounter {
	return &testCounter{
		data: make(map[string]int),
	}
}

func (t *testCounter) Inc(phrase string) {
	t.data[phrase]++
}

func (t *testCounter) GetCounts() map[string]int {
	return t.data
}

func TestCountTop3(t *testing.T) {
	tc := []struct {
		name     string
		r        io.Reader
		c        Counter
		expected map[string]int
		wantErr  bool
	}{
		{
			"will count sections of 3 words",
			bytes.NewBuffer([]byte("123 123 123 321 321 456")),
			newTestCounter(),
			map[string]int{"123 123 123": 1, "123 123 321": 1, "123 321 321": 1, "321 321 456": 1},
			false,
		},
	}

	for _, tc := range tc {
		t.Run(tc.name, func(t *testing.T) {
			err := CountTop3(tc.r, tc.c)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expected, tc.c.GetCounts())
		})
	}
}
