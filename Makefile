.PHONY: check docs fmt test verify

docs:
	test -f docs/plans/2026-06-08-purpleair-go-baseline.md
	grep -q "Status: Completed" docs/plans/2026-06-08-purpleair-go-baseline.md
	grep -q "make check" docs/plans/2026-06-08-purpleair-go-baseline.md

fmt:
	test -z "$$(gofmt -l *.go)"

test:
	go test ./...

verify: fmt test docs

check: verify
