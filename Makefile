.PHONY: check build docs fmt lint race root-test test vet verify

ifneq ($(origin MAKEFILE_LIST),file)
$(error MAKEFILE_LIST must not be overridden)
endif
override REPO_ROOT := $(shell path='$(subst ','"'"',$(MAKEFILE_LIST))'; path=$$(printf '%s' "$$path" | /usr/bin/sed 's/^ //'); directory=$$(/usr/bin/dirname -- "$$path"); CDPATH= cd -- "$$directory" && /bin/pwd -P)

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

root-test:
	cd "$(REPO_ROOT)" && scripts/test-makefile-root.sh

verify: lint vet test race build docs root-test

check: verify
	cd "$(REPO_ROOT)" && scripts/check-baseline.sh
