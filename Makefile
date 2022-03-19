fmt:
	go fmt ./...

vet:
	go vet ./...

test: fmt vet
	go test -race -coverprofile=coverage.out ./...

run: test
	go run main.go

runspecies: test
	go run main.go -file species.txt

benchcmd: fmt vet
	go test -bench=. ./cmd -benchmem -memprofile mem.out -cpuprofile cpu.out

benchscanner: fmt vet
	go test -bench . ./scanner -benchmem -memprofile mem.out -cpuprofile cpu.out

cover: test
	go tool cover -html=coverage.out

concreal:
	go run main.go -concurrent -file huge.txt -file huge.txt -file huge.txt -file huge.txt - file huge_processed.txt -file med_processed.txt -file huge.txt

real:
	go run main.go -file huge.txt -file huge.txt -file huge.txt -file huge.txt - file huge_processed.txt -file med_processed.txt -file huge.txt

build:
	go build -o counttool main.go