.PHONY: build run

build:
	@go build -o bin/api ./cmd/api

run: build            #build is a pre-requisite
	@.\bin\api       

migrate-up:
	@go run ./cmd/migrate up

migrate-down:
	@go run ./cmd/migrate down 