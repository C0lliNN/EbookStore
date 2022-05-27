# E-book Store
![](https://img.shields.io/badge/coverage-95%25-brightgreen)

A fully-featured REST API developed in Golang for an ebook store.

## Table of Contents
* [Features](#features)
* [Tech/Libraries](#techlibraries)
* [High-level Documentation](#high-level-documentation)
  * [System Context Diagram](#system-context-diagram)
  * [Backend Container Diagram](#backend-container-diagram)
  * [Design Principles / Techniques](#design-principles--techniques)
* [Local Development](#local-development)
  * [How to Run Locally](#how-to-run-locally)
  * [How Generate Seed Data](#how-generate-seed-data)
  * [How to Regenerate Mocks](#how-generate-seed-data)
  * [How to Regenerate REST API Documentation](#how-to-regenerate-rest-api-documentation)
  * [How to Regenerate Dependency Initialization](#how-to-regenerate-dependency-initialization)

## Features
* Authentication (Sign up, Login and Reset Password)
* Multiple Roles (Customer and Administrator)
* Book Catalog Management
* Order Management
* Pagination
* Create Orders
* File Storage/Retrieval
* Payment Management 

## Tech/Libraries

* [Golang](https://golang.org/)
* [Gin](https://github.com/gin-gonic/gin)
* [PostgreSQL](https://www.postgresql.org/)
* [GORM](https://gorm.io/index.html)
* [JWT](https://jwt.io/)
* [Bcrypt](https://en.wikipedia.org/wiki/Bcrypt)
* [Wire](https://github.com/google/wire)
* [Viper](https://github.com/spf13/viper)
* [Zap](https://github.com/uber-go/zap)
* [Stripe](https://stripe.com/)
* [Amazon S3](https://aws.amazon.com/s3/?nc1=h_ls)
* [Amazon SES](https://aws.amazon.com/ses/?nc1=h_ls)
* [Swagger](https://www.openapis.org/)
* [Localstack](https://localstack.cloud/)
* [Testify](https://github.com/stretchr/testify)

## High-level Documentation
This is a high-level technical documentation about how this application is structured. The diagrams follow the [C4 model](https://c4model.com/)

### System Context Diagram
![](https://i.ibb.co/Kykm454/TzmXLRz.png)

### Backend Container Diagram
![](https://i.ibb.co/LS9pSDK/image.png)

### Design Principles / Techniques
* SOLID Principles
* Hexagonal Architecture
* Domain Driver Design
* Package by Feature
* Test-Driven-Development

## Local Development

### How to Run Locally
1. Execute docker containers
```bash
make docker-up
```

2. Execute REST HTTP Server
```bash
make start_server
```

3. Open `http://localhost:8080/docs` in your browser

### How Generate Seed Data
1. Execute docker containers
```bash
make docker-up
```

2. Clean database and generate seed data 
```bash
make generate_seed_data
```

### How to Regenerate Mocks
1. Install Mockery
```bash
go install github.com/vektra/mockery/v2@latest
```
2. Generate Mocks
```bash
make generate-mocks
```

### How to Regenerate REST API Documentation
1. Install Swaggo
```bash
go install github.com/swaggo/swag/cmd/swag@v1.7.8
```
2. Generate Docs
```bash
make api-docs
```

### How to Regenerate Dependency Initialization
1. Install wire 
```bash
go install github.com/google/wire/cmd/wire@latest
```
2. Execute wire
```bash
make wire
```