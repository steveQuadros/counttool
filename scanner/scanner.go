package scanner

import (
	"bufio"
	"bytes"
	"io"
)

type Counter interface {
	Inc(string)
	GetCounts() map[string]int
	GetCount(string) int
}

func CountTop3(reader io.Reader, counter Counter) error {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)
	var words int
	var phrase [][]byte
	for scanner.Scan() {
		word := scanner.Bytes()
		words++
		trimWordOfPunctuation(&word)
		word = bytes.ToLower(word)
		phrase = append(phrase, word)
		if words == 3 {
			joined := bytes.Join(phrase, []byte{' '})
			counter.Inc(string(joined))
			phrase = phrase[1:]
			words = 2
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// trimWordOfPunctuation checks if last char is [a-z][A-Z][0-9]
// or removes the last char if not
// a more exhaustive check would need to handle UTF8 sequences, but this should work for a basic case
// TODO - remove punctuation from any point in the word (handle malicious output)
func trimWordOfPunctuation(word *[]byte) {
	lastChar := (*word)[len(*word)-1]
	if !(lastChar >= 'a' && lastChar <= 'z') && !(lastChar >= 'A' && lastChar <= 'Z') && !(lastChar >= '0' && lastChar <= '9') {
		*word = (*word)[:len(*word)-1]
	}
}
