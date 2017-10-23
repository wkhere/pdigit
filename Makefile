go:
	go fmt
	go build
	go vet
	go test -cover
	go install

bench:
	go build
	go test -bench=. -benchmem
	
cover:
	go build
	go test -run=TestCall -coverprofile=cov
	go tool cover -html cov

.PHONY: go bench cover
