go:
	go fmt
	go build # needed because we're also testing Exec
	go test -cover
	go install

bench:
	go build
	go test -bench=. -benchmem
	
cover:
	go test -run=TestCall -coverprofile=cov
	go tool cover -html cov

.PHONY: go bench cover
