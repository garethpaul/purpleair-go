.PHONY: check build docs fmt lint test verify

docs:
	@for plan in docs/plans/*.md; do \
		test -f "$$plan"; \
		grep -q "Status: Completed" "$$plan"; \
		grep -q "make check" "$$plan"; \
	done

fmt:
	test -z "$$(gofmt -l *.go)"

lint: fmt

test:
	go test ./...

build: test

verify: lint test build docs

check: verify
