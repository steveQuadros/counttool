package cmd

import (
	"github.com/stevequadros/counttop/scanner"
	"io"
	"sync"
)

type phraseCount struct {
	data map[string]int
	mu   sync.Mutex
}

func NewPhraseCount() phraseCount {
	return phraseCount{
		data: make(map[string]int),
		mu:   sync.Mutex{},
	}
}

func (p *phraseCount) Inc(phrase string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.data[phrase]++
}

// TODO unsafe to give access to this concurrently
func (p *phraseCount) GetCounts() map[string]int {
	return p.data
}

func Execute(readers []io.ReadCloser) {
	counter := NewPhraseCount()
	for _, r := range readers {
		scanner.CountTop3(r, &counter)
	}
}
