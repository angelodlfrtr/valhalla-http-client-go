GOBIN=go

.PHONY: lint
lint:
	go vet ./...
	golangci-lint run ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: clean
clean:
	rm ./coverage.out

.PHONY: test
test:
	go test -v -count 1 -race --coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
