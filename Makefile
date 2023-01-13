go:
	go vet  ./...
	go test ./...
	go install ./cmd/pdigit

fuzz:
	go test -fuzz=. $(opt)

bench:
	go test -bench=. -benchmem -count=5

cover:
	go test -coverprofile=cov
	go tool cover -html cov

.PHONY: go fuzz bench cover
