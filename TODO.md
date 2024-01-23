# TODO

* [x] Flat config folder
* [x] Place config folder inside internal
* [x] Rename from migration to migrations
* [x] Define Logging Strategy
* [x] Search about ctx wrapper
* [x] Send ctx context.Context in all functions
* [x] Flat auth folder
* [x] Flat catalog folder
* [x] Flat shop folder
* [x] Add validation to usecase layer
* [x] Add Authorization to usecase layer (and add userId to shop layer)
* [x] Double check requests/responses function names looking for inconsistencies
* [x] Create a Route Wrapper
* [x] Create Errors Middleware
* [x] Refactor Errors Middleware
* [x] Create Swagger Middleware
* [x] Make the code compile again / Refactor Wire
* [x] Use pointers type for DI
* [x] Fix Swagger Endpoint
* [x] Place docs folder in root directory
* [x] Remove admin middleware
* [x] Replace `c.Request.Context()` with just `c`
* [x] Refactor migrations folder
* [x] Add health-check
* [x] Review tests (Noticing Naming, Pointer, files that are not tested, check coverage)
* [x] Create a repository testsuite
* [x] Add Logging, Observability
* [x] Review Swagger Documentation
* [x] Review errors (Update Server Error Mapping)
* [x] Reduce number of multi-argument functions
* [x] Review pointer/concrete types in struct and as receivers (tests included)
* [x] Find a way to separate unit and integration tests while keeping good autocompletion
* [x] Document Public API (Maybe autogenerate documentation)
* [x] Try to use context.Context in Stripe
* [x] Test app manually (Save final postman collection somewhere)
* [x] Make test pass (Also remove user_factory)
* [x] Create Seed Script
* [x] Refactor makefile
* [x] Update README.md
* [x] Add golangci-lint
* [x] Add error wrapping
* [x] Implement graceful shutdown
* [x] Change module path to just github.com/ebookstore
* [x] Add Rate Limiting Middleware
* [x] Add Core and Platform folders
* [x] Add Version API
* [x] Add CORS Middleware
* [x] Allow a single a book to have multiple images instead of only one
* [x] Use the S3 Presigning API to send book content to the user instead of downloading the file in the server
* [x] Validate Data Size in All Requests
* [x] Add INFO Logging in important places 
* [x] Protect code against SQL Injection
* [x] Improve test suite, use apitest for the HTTP handlers (Add more catalog tests and add shop tests)
* [x] Close DB in Sever Shutdown
* [x] Organize TODO for next iteration

## Features
* [ ] Create Shopping Cart Functionality
* [ ] Add a position field to the image table
* [ ] Add ability to specify mime type in the request when creating a presigned url

## Refactoring
* [x] Use mocks in the same package
* [x] Refactor Logging (use function log(ctx, msg, fields))

## Architecture Improvements
* [x] Try to implement Query Pattern to decouple domain from Database specific technologies
* [ ] Implement Unit of Work Pattern for transactional use cases
* [ ] Move Request/Response to server package and implement the CQRS Pattern
* [ ] Use one use case per file approach
* [ ] Implement Auth using OAuth2 Specification
* [ ] Use an anti-corruption layer to handle the Stripe Webhook in a better way

## Dev Experience
* [x] Create a new script for cleaning data
* [ ] Use Stripe CLI in local environment (include instructions in the README.md)

## Deploy
* [ ] Add Github Actions capabilities
* [ ] Create prod Dockerfile and deploy folder/script
