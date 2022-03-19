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

func (t *testCounter) GetCount(phrase string) int {
	if count, ok := t.data[phrase]; ok {
		return count
	} else {
		return -1
	}
}

func TestCountTop3(t *testing.T) {
	tc := []struct {
		name            string
		r               io.Reader
		c               Counter
		expected        map[string]int
		wantErr         bool
		testPhrase      *string
		testPhraseCount int
	}{
		{
			"will count sections of overlapping 3 words",
			bytes.NewBuffer([]byte("123 123 123 321 321 456")),
			newTestCounter(),
			map[string]int{"123 123 123": 1, "123 123 321": 1, "123 321 321": 1, "321 321 456": 1},
			false,
			nil,
			0,
		},
		{
			"ignores punctuation at end of words",
			bytes.NewBuffer([]byte("123! 123\n 123, 321 321? 456\n")),
			newTestCounter(),
			map[string]int{"123 123 123": 1, "123 123 321": 1, "123 321 321": 1, "321 321 456": 1},
			false,
			nil,
			0,
		},
		{
			"ignores line endings",
			bytes.NewBuffer([]byte("I love\nsandwiches. I LOVE\nSANDWICHES. I love\nsandwiches.")),
			newTestCounter(),
			map[string]int{"i love sandwiches": 3, "love sandwiches i": 2, "sandwiches i love": 2},
			false,
			strPtr("i love sandwiches"),
			3,
		},
		{
			"is case insensitive",
			bytes.NewBuffer([]byte("A MAN WAKES a Man Wakes a Man wAkeS   \n \n ")),
			newTestCounter(),
			map[string]int{"a man wakes": 3, "man wakes a": 2, "wakes a man": 2},
			false,
			strPtr("a man wakes"),
			3,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			err := CountTop3(tt.r, tt.c)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expected, tt.c.GetCounts())
			if tt.testPhrase != nil {
				require.Equal(t, tt.testPhraseCount, tt.c.GetCount(*tt.testPhrase))
			}
		})
	}
}

func TestTrimLastCharacterPuncutation(t *testing.T) {
	tc := []struct {
		word     []byte
		expected []byte
	}{
		{[]byte("test"), []byte("test")},
		{[]byte("test!"), []byte("test")},
		{[]byte("test."), []byte("test")},
		{[]byte("test?"), []byte("test")},
	}

	for _, tt := range tc {
		t.Run(string(tt.word), func(t *testing.T) {
			trimWordOfPunctuation(&tt.word)
			require.Equal(t, tt.expected, tt.word)
		})
	}
}

func strPtr(s string) *string {
	return &s
}
