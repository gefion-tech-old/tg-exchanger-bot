SOURCES=./cmd/main
SERVICE=bot

.PHONY: init

init:	
	mkdir -p tmp
	mkdir -p logs

.PHONY: run

run:
	go build -o $(SERVICE) -v $(SOURCES)
	clear		
	./$(SERVICE)


.PHONY: build

build:
	go build -o $(SERVICE) -v $(SOURCES)
	clear	


.PHONY: test

test:
	go test -v -race -timeout 30s ./...


count:
	find . -name tests -prune -o -type f -name '*.go' | xargs wc -l

.DEFAULT_GOAL := run