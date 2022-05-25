.PHONY: dependency unit-test test docker-up docker-down wire generate-mocks api-docs start_server

dependency:
	@go get -v ./...

docker-up:
	@docker-compose up -d

docker-down:
	@docker-compose down

start_server: export ENV=local
start_server: export AWS_ACCESS_KEY_ID=test
start_server: export AWS_SECRET_ACCESS_KEY=test
start_server:
	@cd cmd && go run .

test: export ENV=test
test: export AWS_ACCESS_KEY_ID=test
test: export AWS_SECRET_ACCESS_KEY=test
test: dependency
	@go test ./...

test_no_cache: export ENV=test
test_no_cache: export AWS_ACCESS_KEY_ID=test
test_no_cache: export AWS_SECRET_ACCESS_KEY=test
test_no_cache: dependency
	@go test ./... -count=1

unit-test: dependency
	@go test -short ./...

migrate-up:
	@migrate -source file:./migration -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable up

migrate-down:
	@migrate -source file:./migration -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable down 1

wire:
	@cd cmd && wire .

generate-mocks:
	@mockery --all --output=./mocks --dir=internal --case=underscore --keeptree

api-docs:
	@swag init -g internal/server/server.go --parseInternal  --generatedTime
