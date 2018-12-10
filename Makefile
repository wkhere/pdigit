go:
	go fmt
	go test -cover
	go install

bench:
	go test -bench=. -benchmem

cover:
	go test -coverprofile=cov
	go tool cover -html cov

.PHONY: go bench cover
