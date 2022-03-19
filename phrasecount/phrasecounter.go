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

type Count struct {
	data map[string]int
}

var _ Counter = (*Count)(nil)

func NewPhraseCount() Count {
	return Count{
		data: make(map[string]int),
	}
}

func (p *Count) Inc(phrase string) {
	p.data[phrase]++
}

func (p *Count) GetCounts() map[string]int {
	return p.data
}

func (p *Count) GetCount(phrase string) int {
	if count, ok := p.data[phrase]; ok {
		return count
	} else {
		return -1
	}
}

func (p *Count) Top(n int) PhraseOutputList {
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

type CountConcurrent struct {
	data map[string]int
	mu   sync.RWMutex
}

var _ Counter = (*CountConcurrent)(nil)

func NewPhraseCountConcurrent() CountConcurrent {
	return CountConcurrent{
		data: make(map[string]int),
		mu:   sync.RWMutex{},
	}
}

func (p *CountConcurrent) Inc(phrase string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.data[phrase]++
}

// TODO unsafe to give access to this concurrently
func (p *CountConcurrent) GetCounts() map[string]int {
	return p.data
}

func (p *CountConcurrent) GetCount(phrase string) int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if count, ok := p.data[phrase]; ok {
		return count
	} else {
		return -1
	}
}

func (p *CountConcurrent) Top(n int) PhraseOutputList {
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
