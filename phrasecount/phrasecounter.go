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
}

var _ Counter = (*PhraseCount)(nil)

func NewPhraseCount() PhraseCount {
	return PhraseCount{
		data: make(map[string]int),
	}
}

func (p *PhraseCount) Inc(phrase string) {
	p.data[phrase]++
}

func (p *PhraseCount) GetCounts() map[string]int {
	return p.data
}

func (p *PhraseCount) GetCount(phrase string) int {
	if count, ok := p.data[phrase]; ok {
		return count
	} else {
		return -1
	}
}

func (p *PhraseCount) Top(n int) PhraseOutputList {
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

type PhraseCountConcurrent struct {
	data map[string]int
	mu   sync.RWMutex
}

var _ Counter = (*PhraseCountConcurrent)(nil)

func NewPhraseCountConcurrent() PhraseCountConcurrent {
	return PhraseCountConcurrent{
		data: make(map[string]int),
		mu:   sync.RWMutex{},
	}
}

func (p *PhraseCountConcurrent) Inc(phrase string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.data[phrase]++
}

// TODO unsafe to give access to this concurrently
func (p *PhraseCountConcurrent) GetCounts() map[string]int {
	return p.data
}

func (p *PhraseCountConcurrent) GetCount(phrase string) int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if count, ok := p.data[phrase]; ok {
		return count
	} else {
		return -1
	}
}

func (p *PhraseCountConcurrent) Top(n int) PhraseOutputList {
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
