.PHONY: run automigrate dev infra test

run:
	@go run cmd/server/main.go

automigrate:
	@go run cmd/automigrate/main.go config.json

dev:
	air

infra:
	docker-compose up -d

test:
	@go test -v -race -coverprofile /tmp/c.out ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''
