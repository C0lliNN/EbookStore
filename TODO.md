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
* [ ] Fix Swagger Endpoint
* [ ] Place docs folder in root directory
* [ ] Refactor migrations folder
* [ ] Replace `c.Request.Context()` with just `c`
* [ ] Remove admin middleware
* [ ] Review tests (Noticing Naming, Pointer, files that are not tested, check coverage)
* [ ] Add Logging, Observability
* [ ] Review Swagger Documentation
* [ ] Review errors (Update Server Error Mapping)
* [ ] Reduce number of multi-argument functions
* [ ] Review pointer/concrete types in struct and as receivers (tests included)
* [ ] Find a way to separate unit and integration tests while keeping good autocompletion
* [ ] Create Good Seed Data
* [ ] Document Public API (Maybe autogenerate documentation)
* [ ] Get rid of test folder
* [ ] Add oauth server
* [ ] Refactor make command
* [ ] Create make run command and update README.md

# Nice to Haves
* Handle Stripe Webhook in a good way
* Handle file extensions in a more isolated way
* Handle shop -> catalog connection isolating the BookResponse object
* Use testcontainers in integration tests

(Update Handler tests)