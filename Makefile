.PHONY: test one fmt bench fetch build

PKG := $(subst -,_,$(name))

# Run all tests
test:
	go test ./... -v

# Run tests for a single problem: make one name=two_sum
one:
	@test -n "$(name)" || (echo "usage: make one name=two_sum" && exit 1)
	go test ./problems/$(name)/... -v

# Fetch a problem from LeetCode and scaffold it.
# Usage:
#   make fetch url=https://leetcode.com/problems/valid-parentheses/
#   make fetch slug=valid-parentheses
fetch: build
	@test -n "$(url)$(slug)" || (echo "usage: make fetch url=https://leetcode.com/problems/... OR make fetch slug=two-sum" && exit 1)
	@if [ -n "$(url)" ]; then \
		./bin/fetch -url "$(url)"; \
	else \
		./bin/fetch -slug "$(slug)"; \
	fi

# Build the fetch binary
build:
	@mkdir -p bin
	go build -o bin/fetch ./cmd/fetch

# Format all Go files
fmt:
	gofmt -w $$(find . -type f -name '*.go' -not -path './vendor/*')

# Run benchmarks
bench:
	go test ./... -bench=. -benchmem
