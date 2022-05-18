go:
	go fmt
	go test ./...
	go install

fuzz:
	go test -fuzz=. $(opt)

bench:
	go test -bench=. -benchmem

cover:
	go test -coverprofile=cov
	go tool cover -html cov

.PHONY: go fuzz bench cover
