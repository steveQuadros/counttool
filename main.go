package main

import (
	"flag"
	"fmt"
	"github.com/stevequadros/counttop/cmd"
	"io"
	"os"
	"strings"
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
	flag.Var(&filesFlag, "file", "Files to count phrases in. usage: ./counttop --file file1.txt --file file2.txt ")
	flag.IntVar(&topCount, "count", 100, "number of top results to return. Default to 3.")
	flag.Parse()

	var in []io.ReadCloser
	if len(filesFlag) == 0 {
		// assume reading from stdin
		in = append(in, os.Stdin)
	} else {
		// assume reading from files
		for _, f := range filesFlag {
			file, err := os.Open(f)
			if err != nil {
				logErrorAndExit("error opening file", err)
			}
			in = append(in, file)
		}
	}

	top, err := cmd.Execute(in, topCount)
	if err != nil {
		logErrorAndExit("error parsing files", err)
	}
	for _, res := range top {
		logSuccess(fmt.Sprintf("%d - %q", res.Count, res.Phrase))
	}
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
