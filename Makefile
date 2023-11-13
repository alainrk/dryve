.PHONY: test automigrate start-db dev start

automigrate:
	@go run cmd/automigrate/main.go

start-db:
	docker-compose up -d db

dev:
	air

test:
	@go test -v -race -coverprofile /tmp/c.out ./...

start:
	./setup.sh

run:
	go run cmd/server/main.go
