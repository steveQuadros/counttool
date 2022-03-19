# REQUIREMENTS
- [x] The program accepts as arguments a list of one or more file paths (e.g. ./solution.rb file1.txt
file2.txt ...).
- [x] The program also accepts input on stdin (e.g. cat file1.txt | ./solution.rb).

- [x] The program outputs a list of the 100 most common three word sequences in the text, along
with a count of how many times each occurred in the text. For example: `231 - i will not, 116 - i
do not, 105 - there is no, 54 - i know not, 37 - i am not ...`
- [x] The program ignores punctuation, line endings, and is case insensitive (e.g. “I
love\nsandwiches.” should be treated the same as "(I LOVE SANDWICHES!!)")
- [x] The program is capable of processing large files and runs as fast as possible.
- [x] The program should be tested. Provide a test file for your solution.
- [x] The program should be well structured and understandable.

## QuickStart
`make build`

`./counttool [-concurrent] [-count] [-file...] [-file]`

ex.
`./counttool -concurrent -file huge.txt -file huge.txt -file huge.txt -file huge.txt - file huge_processed.txt -file med_processed.txt -file huge.txt `

or 

`cat huge.txt | ./counttool [-concurrent] [-count]`

I also added the ability to pre-process a file to see if preprocessing would significantly impact speed (Narrator: it did not)
Preprocessing downcases all words and removes all puncuation, leaving only text and single spaces
This will output to stdout, so you can redirect where you'd like

`./counttool -process -file small.txt > small_processed.txt`

## Benchmarks

Interestingly concurrent reading / processing files benchmarks more slowly
```
BenchmarkExecute-8             	  327913	      3485 ns/op	   16653 B/op	       9 allocs/op
BenchmarkExecuteConcurrent-8   	  145083	      8099 ns/op	   17042 B/op	      15 allocs/op
```

Whereas real world tests against the huge files show it being faster when processing many large files. To see this run:

`make real` for non-concurrent

then 

`make concreal` for concurrent
