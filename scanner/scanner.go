package scanner

import (
	"bufio"
	"bytes"
	"io"
)

type Counter interface {
	Inc(string)
	GetCounts() map[string]int
}

func CountTop3(reader io.Reader, counter Counter) error {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)
	var words int
	var phrase [][]byte
	for scanner.Scan() {
		word := scanner.Bytes()
		words++
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
