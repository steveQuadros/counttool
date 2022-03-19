package cmd

import (
	"github.com/stevequadros/counttop/phrasecount"
	"github.com/stevequadros/counttop/scanner"
	"io"
)

func Execute(readers []io.Reader, count int) (phrasecount.PhraseOutputList, error) {
	counter := phrasecount.NewPhraseCount()
	for _, r := range readers {
		if err := func() error {
			if err := scanner.CountTop(r, &counter); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			return nil, err
		}
	}
	return counter.Top(count), nil
}

func ExecuteConcurrent(readers []io.Reader, count, workers int) (phrasecount.PhraseOutputList, error) {
	counter := phrasecount.NewPhraseCountConcurrent()
	errors := make(chan error, len(readers))
	defer close(errors)
	for _, r := range readers {
		go func(reader io.Reader) {
			errors <- scanner.CountTop(reader, &counter)
		}(r)
	}

	var doneCount int
	for doneCount < len(readers) {
		select {
		case err := <-errors:
			doneCount++
			if err != nil {
				return nil, err
			}
		}
	}
	return counter.Top(count), nil
}
