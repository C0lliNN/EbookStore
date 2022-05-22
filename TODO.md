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
* [ ] Add Logging, Observability
* [x] Review Swagger Documentation
* [x] Review errors (Update Server Error Mapping)
* [ ] Reduce number of multi-argument functions
* [ ] Review pointer/concrete types in struct and as receivers (tests included)
* [x] Find a way to separate unit and integration tests while keeping good autocompletion
* [ ] Create Good Seed Data
* [x] Document Public API (Maybe autogenerate documentation)
* [ ] Get rid of test folder
* [ ] Add oauth server
* [ ] Add more integration tests
* [ ] Refactor make command
* [ ] Create make run command and update README.md

# Nice to Haves
* Handle Stripe Webhook in a good way
* Handle file extensions in a more isolated way
* Handle shop -> catalog connection isolating the BookResponse object
* Handle more than one image per book
* Use Cart in Shop 
(Update Handler tests)