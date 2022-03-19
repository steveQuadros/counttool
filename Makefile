fmt:
	go fmt ./...

vet:
	go vet ./...

test: fmt vet
	go test -race -cover ./...

run: test
	go run main.go

runspecies: test
	go run main.go -file species.txt

bench: fmt vet
	go test -bench=. ./cmd -benchmem -memprofile mem.out -cpuprofile cpu.out