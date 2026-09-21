.PHONY: build run

build:
	@go build -o bin/api ./cmd/api

run: build            #build is a pre-requisite
	@.\bin\api       