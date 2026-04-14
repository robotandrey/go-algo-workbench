.PHONY: test one fmt bench

test:
	go test ./... -v

one:
	@test -n "$(name)" || (echo "usage: make one name=two_sum" && exit 1)
	go test ./problems/$(name) -v

fmt:
	gofmt -w $$(find . -type f -name '*.go' -not -path './vendor/*')

bench:
	go test ./... -bench=. -benchmem
