package cmd

import (
	"github.com/stevequadros/counttop/phrasecount"
	"github.com/stevequadros/counttop/scanner"
	"io"
)

func Execute(readers []io.ReadCloser, count int) ([]phrasecount.PhraseOutput, error) {
	counter := phrasecount.NewPhraseCount()
	for _, r := range readers {
		// wrap so defer statement closes file as quickly as possible, see: https://stackoverflow.com/a/45620423
		if err := func() error {
			if err := scanner.CountTop(r, &counter); err != nil {
				return err
			}
			defer r.Close()
			return nil
		}(); err != nil {
			return nil, err
		}
	}
	return counter.Top(count), nil
}
