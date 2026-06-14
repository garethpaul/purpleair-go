.PHONY: check build docs fmt lint race test vet verify

override REPO_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

docs:
	@cd "$(REPO_ROOT)" && for plan in docs/plans/*.md; do \
		test -f "$$plan"; \
		grep -q "Status: Completed" "$$plan"; \
		grep -q "make check" "$$plan"; \
	done

fmt:
	cd "$(REPO_ROOT)" && test -z "$$(gofmt -l *.go)"

lint: fmt

vet:
	cd "$(REPO_ROOT)" && go vet ./...

test:
	cd "$(REPO_ROOT)" && go test ./...

race:
	cd "$(REPO_ROOT)" && go test -race ./...

build: test

verify: lint vet test race build docs

check: verify
	cd "$(REPO_ROOT)" && scripts/check-baseline.sh
