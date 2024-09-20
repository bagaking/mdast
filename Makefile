.PHONY: check fmt test vet check-imports

check: fmt test vet check-imports

fmt:
	test -z "$$(gofmt -l .)"

test:
	go test ./...

vet:
	go vet ./...

check-imports:
	@imports="$$(go list -f '{{range .Imports}}{{println .}}{{end}}' .)"; \
	if printf '%s\n' "$$imports" | grep -E '^(testing|github.com/google/go-cmp|github.com/stretchr/testify)($$|/)'; then \
		echo "production package imports testing or assertion libraries"; \
		exit 1; \
	fi
