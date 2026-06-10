.PHONY: check build docs fmt lint race test vet verify

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

race:
	go test -race ./...

build: test

verify: lint vet test race build docs

check: verify
	scripts/check-baseline.sh
