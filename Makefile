.PHONY: test

test:
	go test -v ./... --cover
	go vet ./...
