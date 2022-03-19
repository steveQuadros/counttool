package scanner

import (
	"bufio"
	"bytes"
	"github.com/stevequadros/counttop/phrasecount"
	"io"
	"unicode"
	"unicode/utf8"
)

func CountTop(reader io.Reader, counter phrasecount.Counter) error {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)
	var phrase [][]byte
	for scanner.Scan() {
		word := scanner.Bytes()
		downcased := bytes.ToLower(word)
		noPunc := RemovePunc(downcased)
		if len(noPunc) != 0 {
			phrase = append(phrase, noPunc)
		}

		if len(phrase) == 3 {
			joined := bytes.Join(phrase, []byte{' '})
			counter.Inc(string(joined))
			phrase = phrase[1:]
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// RemovePunc removes punctuation from any words
func RemovePunc(word []byte) []byte {
	runes := bytes.Runes(word)
	bs := make([]byte, len(runes)*utf8.UTFMax)

	count := 0
	for _, r := range runes {
		if !unicode.IsPunct(r) {
			count += utf8.EncodeRune(bs[count:], r)
		}
	}
	bs = bs[:count]
	return bs
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
