.PHONY: fmt test verify

fmt:
	test -z "$$(gofmt -l *.go)"

test:
	go test ./...

verify: fmt test
