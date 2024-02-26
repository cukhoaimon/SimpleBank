# SimpleBank

## Architecture
<img width="760" alt="image" src="https://github.com/cukhoaimon/SimpleBank/assets/60815035/d29a0982-cc5b-429b-a477-2518fd2095b2">

 
# For local development
## Requirements
- [gRPC](https://grpc.io): create gRPC api
- [gRPC-Gateway](https://grpc-ecosystem.github.io/grpc-gateway/): serve HTTP request along with gRPC request.
- [Protobuf](https://protobuf.dev): define gRPC api contract
- [Docker](https://www.docker.com): run database
- [GoMock](https://github.com/golang/mock): create test db
- [Paseto](https://github.com/o1egl/paseto): generate paseto token
- [dbdocs](https://dbdocs.io): create documentation for database base on database.dbml file.
- [swagger](https://swagger.io): create documentation for api. 

## Documentation
### Database
- Run in terminal:
```
make dbdocs
```
- Follow the guide and click the link showed on terminal.

### Database
- Run in terminal:
```
make server
```
- The default link to documentation is `localhost:8080/swagger/index.html`

## Makefile
- Run database container:
```
make postgres
```
- Run migration up:
```
make migrate-up
```
- Run migration down:
```
make migrate-down
```
- Generate Go code from *.sql file located in ./db/query/*.sql
```
make sqlc
```
- Generate database documentation
```
make dbdocs
```
- Generate mock db
```
make mock
```
- Generate go code and documentation for api from *.proto file located in ./proto/*.proto
```
make proto
```
