package cmd

import (
	"github.com/stevequadros/counttop/phrasecount"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name    string
		readers []io.Reader
		count   int
		want    phrasecount.PhraseOutputList
		wantErr bool
	}{
		{
			"basic execute",
			getFiles(t),
			3,
			phrasecount.PhraseOutputList{
				{Phrase: "of the same", Count: 333},
				{Phrase: "conditions of life", Count: 134},
				{Phrase: "the same species", Count: 128},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Execute(tt.readers, tt.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.ElementsMatch(t, got, tt.want)
		})
	}
}

func TestExecuteConcurrent(t *testing.T) {
	tests := []struct {
		name    string
		readers []io.Reader
		count   int
		want    phrasecount.PhraseOutputList
		wantErr bool
	}{
		{
			"basic execute",
			getFiles(t),
			3,
			phrasecount.PhraseOutputList{
				{Phrase: "of the same", Count: 333},
				{Phrase: "conditions of life", Count: 134},
				{Phrase: "the same species", Count: 128},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExecuteConcurrent(tt.readers, tt.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.ElementsMatch(t, got, tt.want)
		})
	}
}

func BenchmarkExecute(b *testing.B) {
	in := getFiles(b)
	for i := 0; i < b.N; i++ {
		_, err := Execute(in, 100)
		require.NoError(b, err)
	}
}

func BenchmarkExecuteConcurrent(b *testing.B) {
	in := getFiles(b)
	for i := 0; i < b.N; i++ {
		_, err := ExecuteConcurrent(in, 100)
		require.NoError(b, err)
	}
}

type testBT interface {
	FailNow()
	Errorf(string, ...interface{})
}

func getFiles(t testBT) []io.Reader {
	smallEx, err := os.Open("../small.txt")
	require.NoError(t, err)
	medEx, err := os.Open("../med.txt")
	require.NoError(t, err)
	largeEx, err := os.Open("../large.txt")
	require.NoError(t, err)
	fullEx, err := os.Open("../species.txt")
	require.NoError(t, err)
	return []io.Reader{smallEx, medEx, largeEx, fullEx}
}
