.PHONY: dependency unit-test integration-test docker-up docker-down clear

dependency:
	@go get -v ./...

integration-test: docker-up dependency
	@go test -tags=integration ./...

unit-test: dependency
	@go test -tags=unit ./...

docker-up:
	@docker-compose --file=docker-compose.test.yml up -d

docker-down:
	@docker-compose --file=docker-compose.test.yml down

clear: docker-down