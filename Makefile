all: gollum

.PHONY: bin/golangci-lint
bin/golangci-lint: 
	GOBIN=`pwd`/bin go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

.PHONY: bin/gofumpt
bin/gofumpt:
	GOBIN=`pwd`/bin go install mvdan.cc/gofumpt@latest

.PHONY: bin
bin: bin/golangci-lint bin/gofumpt

.PHONY: lint
lint: 
	./bin/golangci-lint run 

.PHONY: fmt
fmt: 
	./bin/gofumpt -l .

.PHONY: gollum
gollum:
	go build .

.PHONY: clean
clean: 
	rm gollum