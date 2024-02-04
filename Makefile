postgres: 
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

createdb: 
	docker exec -it postgres16 createdb --username=root --owner=root simple_bank

dropdb: 
	docker exec -it postgres16 dropdb simple_bank

migrate-up:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrate-up1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migrate-down1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc: 
	sqlc generate

test: 
	go test -v -cover ./...

server:
	go run main.go

dbdocs:
	dbdocs build doc/database.dbml

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/cukhoaimon/SimpleBank/db/sqlc Store

.PHONY: postgres createdb dropdb migrate-up migrate-down sqlc test server mock dbdocs