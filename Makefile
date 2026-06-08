.PHONY: check docs fmt test verify

docs:
	@for plan in docs/plans/*.md; do \
		test -f "$$plan"; \
		grep -q "Status: Completed" "$$plan"; \
		grep -q "make check" "$$plan"; \
	done

fmt:
	test -z "$$(gofmt -l *.go)"

test:
	go test ./...

verify: fmt test docs

check: verify
