go:
	go fmt
	go build
	go vet
	go test
	go install

bench:
	go build
	go test -bench=. -benchmem
	

.PHONY: go bench
