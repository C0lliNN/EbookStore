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
* [ ] Add golangci-lint
* [ ] Create prod Dockerfile and deploy folder/script
* [ ] Add error wrapping
* [ ] Implement graceful shutdown
* [ ] Change module path to just github.com/ebookstore
* [ ] Consider Moving DTO to the HTTP Handlers (Naming them Requests and Use Commands in the UseCase Layer)
* [ ] Consider Using one file per usecase
* [ ] Version API
* [ ] Rate Limiting
* [ ] Review SQL Injection Concerns
* [ ] Review CORS
* [ ] Validate Data Size in All Requests
* [ ] Review Logging, Monitoring

# Nice to Haves
* Handle Stripe Webhook in a good way
* Handle file extensions in a more isolated way
* Handle shop -> catalog connection isolating the BookResponse object
* Handle more than one image per book (Use Presigning instead of receiving images in the request body)
* Use Cart in Shop 
* Add oauth server
* Add more integration tests
* Use Stripe-CLI to the local environment (including a how to update payment intent in the cli)
* Make Stripe Webhook endpoint safe by querying Stripe API using the provided information instead of asserting it's true
* Use Specification in Query Commands
* Improve mocks location
* Create Auth, Catalog and Shop Diagrams
