.PHONY: check build docs fmt lint test vet verify

docs:
	@for plan in docs/plans/*.md; do \
		test -f "$$plan"; \
		grep -q "Status: Completed" "$$plan"; \
		grep -q "make check" "$$plan"; \
	done

fmt:
	test -z "$$(gofmt -l *.go)"

lint: fmt

vet:
	go vet ./...

test:
	go test ./...

build: test

verify: lint vet test build docs

check: verify
	scripts/check-baseline.sh
