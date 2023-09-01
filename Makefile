target = ./cmd/pdigit

default: test

build:
	go build $(target)

install: test
	go install $(target)

test:
	go test ./... $(opt)

fuzz:
	go test -fuzz=. $(opt)

bench:
	go test -bench=$(sel) -benchmem -count=$(cnt)
sel=.
cnt=5

cover:
	go test -coverprofile=cov
	go tool cover -html=cov -o cov.html && browse cov.html

.PHONY: build install test fuzz bench cover
