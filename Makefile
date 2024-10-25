all: gollum

.PHONY: bin/staticcheck
bin/staticcheck: 
	GOBIN=`pwd`/bin go install honnef.co/go/tools/cmd/staticcheck@latest

.PHONY: bin
bin: bin/staticcheck

.PHONY: gollum
gollum:
	go build .

.PHONY: clean
clean: 
	rm gollum