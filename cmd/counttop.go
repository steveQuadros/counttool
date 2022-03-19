package cmd

import (
	"github.com/stevequadros/counttop/scanner"
	"io"
	"sort"
	"sync"
)

type PhraseCount struct {
	data map[string]int
	mu   sync.RWMutex
}

func NewPhraseCount() PhraseCount {
	return PhraseCount{
		data: make(map[string]int),
		mu:   sync.RWMutex{},
	}
}

func (p *PhraseCount) Inc(phrase string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.data[phrase]++
}

// TODO unsafe to give access to this concurrently
func (p *PhraseCount) GetCounts() map[string]int {
	return p.data
}

func (p *PhraseCount) GetCount(phrase string) int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if count, ok := p.data[phrase]; ok {
		return count
	} else {
		return -1
	}
}

type PhraseOutput struct {
	Phrase string
	Count  int
}

func (p *PhraseCount) Top(n int) []PhraseOutput {
	p.mu.Lock()
	defer p.mu.Unlock()
	var all []PhraseOutput
	for phrase, count := range p.data {
		all = append(all, PhraseOutput{phrase, count})
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].Count > all[j].Count
	})

	if n >= len(all) {
		return all
	} else {
		return all[:n]
	}
}

func Execute(readers []io.ReadCloser, count int) ([]PhraseOutput, error) {
	counter := NewPhraseCount()
	for _, r := range readers {
		// wrap so defer statement closes file as quickly as possible, see: https://stackoverflow.com/a/45620423
		if err := func() error {
			if err := scanner.CountTop3(r, &counter); err != nil {
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
