package phrasecount

import (
	"sort"
	"sync"
)

type Counter interface {
	Inc(string)
	GetCounts() map[string]int
	GetCount(string) int
	Top(int) PhraseOutputList
}

type PhraseCount struct {
	data map[string]int
	mu   sync.RWMutex
}

var _ Counter = (*PhraseCount)(nil)

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

type PhraseOutputList []PhraseOutput

func (p PhraseOutputList) Map() map[string]int {
	m := make(map[string]int)
	for _, po := range p {
		m[po.Phrase] = po.Count
	}
	return m
}

type PhraseOutput struct {
	Phrase string
	Count  int
}

func (p *PhraseCount) Top(n int) PhraseOutputList {
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
