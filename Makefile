target = ./cmd/pdigit

default: test build

build:
	go build $(target)

install: test
	go install $(target)

test:
	go test ./... $(opt)

fuzz:
	go test -fuzz=. $(opt)

sel=.
cnt=5

bench:
	go test -bench=$(sel) -benchmem -count=$(cnt)

cover:
	go test -coverprofile=cov -run $(sel)
	go tool cover -html=cov -o cov.html && browse cov.html

.PHONY: build install test fuzz bench cover
