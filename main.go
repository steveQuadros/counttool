package main

import (
	"flag"
	"fmt"
	"github.com/stevequadros/counttop/cmd"
	"github.com/stevequadros/counttop/phrasecount"
	"io"
	"os"
	"strings"
	"time"
)

type FilesFlag []string

func (f *FilesFlag) Set(s string) error {
	*f = append(*f, s)
	return nil
}

func (f *FilesFlag) String() string {
	b := strings.Builder{}
	for _, p := range *f {
		b.WriteString(p)
	}
	return b.String()
}

func main() {
	var filesFlag FilesFlag
	var topCount int
	var workers int
	flag.Var(&filesFlag, "file", "Files to count phrases in. usage: ./counttop --file file1.txt --file file2.txt ")
	flag.IntVar(&topCount, "count", 100, "number of top results to return. Default to 3.")
	flag.IntVar(&workers, "workers", 1, "number of go routines to use")
	flag.Parse()

	start := time.Now()

	var in []io.Reader
	var fileCount int
	if len(filesFlag) == 0 {
		fileCount = 1
		// assume reading from stdin
		in = append(in, os.Stdin)
	} else {
		fileCount = len(filesFlag)
		// assume reading from files
		for _, f := range filesFlag {
			file, err := os.Open(f)
			if err != nil {
				logErrorAndExit("error opening file", err)
			}
			in = append(in, file)
		}
	}

	var top phrasecount.PhraseOutputList
	var err error
	if workers == 1 {
		top, err = cmd.Execute(in, topCount)
	} else {
		top, err = cmd.ExecuteConcurrent(in, topCount, workers)
	}

	if err != nil {
		logErrorAndExit("error parsing files", err)
	}
	for _, res := range top {
		logSuccess(fmt.Sprintf("%d - %q", res.Count, res.Phrase))
	}

	fmt.Printf("Processed %d files in %.04fs\n", fileCount, time.Now().Sub(start).Seconds())
	os.Exit(0)
}

func logInProcess(s string) {
	fmt.Println(s + "...")
}

func logError(prepend string, e error) {
	fmt.Println("\t✗ " + prepend)
	fmt.Println("\t\t", e)
}

func logErrorAndExit(prepend string, e error) {
	logError(prepend, e)
	os.Exit(1)
}

func logSuccess(s ...string) {
	for _, t := range s {
		fmt.Println("\t✓ " + t)
	}
}
