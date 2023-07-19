go:
	go build ./cmd/pdigit
	go test ./...

install:
	go install ./cmd/pdigit

fuzz:
	go test -fuzz=. $(opt)

bench:
	go test -bench=$(sel) -benchmem -count=$(cnt)
sel=.
cnt=5

cover:
	go test -coverprofile=cov
	go tool cover -html cov

.PHONY: go install fuzz bench cover
