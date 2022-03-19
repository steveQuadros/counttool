package scanner

import (
	"bytes"
	"github.com/stevequadros/counttop/phrasecount"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"sort"
	"testing"
)

type testCounter struct {
	data map[string]int
}

var _ phrasecount.Counter = (*testCounter)(nil)

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

func (t *testCounter) Top(n int) phrasecount.PhraseOutputList {
	var all []phrasecount.PhraseOutput
	for p, c := range t.data {
		all = append(all, phrasecount.PhraseOutput{Phrase: p, Count: c})
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Count > all[j].Count })
	if n == -1 || n >= len(t.data) {
		return all
	} else {
		return all[:n]
	}
}

func TestCountTop3(t *testing.T) {
	smallEx, err := os.Open("../small.txt")
	require.NoError(t, err)
	medEx, err := os.Open("../med.txt")
	require.NoError(t, err)
	largeEx, err := os.Open("../large.txt")
	require.NoError(t, err)
	fullEx, err := os.Open("../species.txt")

	tc := []struct {
		name            string
		r               io.Reader
		c               phrasecount.Counter
		expected        map[string]int
		wantErr         bool
		testPhrase      *string
		testPhraseCount int
		topCount        int
	}{
		{
			"will count sections of overlapping 3 words",
			bytes.NewBuffer([]byte("123 123 123 321 321 456 ")),
			newTestCounter(),
			map[string]int{"123 123 123": 1, "123 123 321": 1, "123 321 321": 1, "321 321 456": 1},
			false,
			nil,
			0,
			-1,
		},
		{
			"ignores punctuation at end of words",
			bytes.NewBuffer([]byte("123! 123\n 123, 321 321? 456\n")),
			newTestCounter(),
			map[string]int{"123 123 123": 1, "123 123 321": 1, "123 321 321": 1, "321 321 456": 1},
			false,
			nil,
			0,
			-1,
		},
		{
			"ignores line endings",
			bytes.NewBuffer([]byte("I love\nsandwiches. I LOVE\nSANDWICHES. I love\nsandwiches.")),
			newTestCounter(),
			map[string]int{"i love sandwiches": 3, "love sandwiches i": 2, "sandwiches i love": 2},
			false,
			strPtr("i love sandwiches"),
			3,
			-1,
		},
		{
			"is case insensitive",
			bytes.NewBuffer([]byte("A MAN WAKES a Man Wakes a Man wAkeS   \n \n ")),
			newTestCounter(),
			map[string]int{"a man wakes": 3, "man wakes a": 2, "wakes a man": 2},
			false,
			strPtr("a man wakes"),
			3,
			-1,
		},
		{
			"small example",
			smallEx,
			newTestCounter(),
			map[string]int{"on the origin": 3, "origin of species": 3, "the origin of": 3, "the project gutenberg": 4},
			false,
			strPtr("the project gutenberg"),
			4,
			4,
		},
		{
			"med example",
			medEx,
			newTestCounter(),
			map[string]int{"on the origin": 6, "origin of species": 6, "the origin of": 6, "the project gutenberg": 4},
			false,
			strPtr("the project gutenberg"),
			4,
			4,
		},
		{
			"large example",
			largeEx,
			newTestCounter(),
			map[string]int{"of natural selection": 15, "on the origin": 14, "origin of species": 16, "the origin of": 17},
			false,
			strPtr("origin of species"),
			16,
			4,
		},
		{
			"full example",
			fullEx,
			newTestCounter(),
			map[string]int{"conditions of life": 125, "in the same": 116, "of the same": 320, "the same species": 126},
			false,
			strPtr("of the same"),
			320,
			4,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			err := CountTop(tt.r, tt.c)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expected, tt.c.Top(tt.topCount).Map())
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
