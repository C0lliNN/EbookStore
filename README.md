# E-book Store
![](https://img.shields.io/badge/coverage-95%25-brightgreen)

A fully-featured REST API developed in Golang for an online book store.

##### [API Documentation](https://ebook-store2.herokuapp.com/swagger/index.html)

## Features
* Authentication (Sign up, Login and Reset Password)
* Multiple Roles (Customer and Administrator)
* Pagination
* Create book
* List books
* Create Orders
* Download books

## Tech/Libraries

* [Golang](https://golang.org/)
* [Gin](https://github.com/gin-gonic/gin)
* [PostgreSQL](https://www.postgresql.org/)
* [Swagger](https://www.openapis.org/)
* [JWT](https://jwt.io/)
* [Bcrypt](https://en.wikipedia.org/wiki/Bcrypt)
* [Wire](https://github.com/google/wire)
* [Viper](https://github.com/spf13/viper)
* [Stripe](https://stripe.com/)
* [Amazon S3](https://aws.amazon.com/s3/?nc1=h_ls)
* [Amazon SES](https://aws.amazon.com/ses/?nc1=h_ls)
* [Localstack](https://localstack.cloud/)
* [Testify](https://github.com/stretchr/testify)

## How to run
1. Execute docker containers
```bash
docker-compose up -d
```
2. In order to execute AWS S3 locally, it is necessary to export the following environment variables:
```bash
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
```

3. Execute app
```bash
cd cmd && go run .
```

## How to regenerate mocks
1. Install Mockery
```bash
go install github.com/vektra/mockery/v2@latest
```
2. Generate Mocks
```bash
make generate-mocks
```

## How to regenerate REST API documentation
1. Install Swaggo
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```
2. Generate Docs
```bash
make api-docs
```

## How to recompile dependencies
1. Install wire 
```bash
go install github.com/google/wire/cmd/wire@latest
```
2. Execute wire
```bash
make wire
```