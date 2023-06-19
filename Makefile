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


lint:
	golangci-lint run

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
	@swag init -g internal/platform/server/server.go --parseInternal  --generatedTime

generate_seed_data: export ENV=local
generate_seed_data: export AWS_ACCESS_KEY_ID=test
generate_seed_data: export AWS_SECRET_ACCESS_KEY=test
generate_seed_data: export DATABASE_URI=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
generate_seed_data: export MIGRATION_SOURCE=file:../../migrations
generate_seed_data: export AWS_REGION=us-east-2
generate_seed_data: export AWS_S3_BUCKET=ebook-store
generate_seed_data: export AWS_S3_ENDPOINT=http://s3.localhost.localstack.cloud:4566
generate_seed_data: export STRIPE_API_KEY=sk_test_51HAKIGHKmAtjDhlfifsr2lIoY8nQZXkQTE2RvqFfa4ASe6Rlk4YRfVxp44Rr9eeSrPivk55dloy9KFv5Zal3sWQz009q9hiu1u
generate_seed_data: dependency
	@cd scripts/generate_seed_data && go run .
