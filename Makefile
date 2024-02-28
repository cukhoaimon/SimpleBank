postgres: 
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

createdb: 
	docker exec -it postgres16 createdb --username=root --owner=root simple_bank

dropdb: 
	docker exec -it postgres16 dropdb simple_bank

migrate-up:
	migrate -path migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrate-down:
	migrate -path migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrate-up1:
	migrate -path migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migrate-down1:
	migrate -path migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc: 
	sqlc generate

test: 
	go test -v -cover ./...

server:
	go run cmd/main.go

dbdocs:
	dbdocs build docs/database.dbml

mock:
	mockgen -package mockdb -destination internal/delivery/http/mock/store.go github.com/cukhoaimon/SimpleBank/internal/usecase/sqlc Store

proto:
	del docs\swagger\*.swagger.json &
	del internal\delivery\grpc\pb &
	protoc --proto_path=internal/delivery/grpc/proto --go_out=internal/delivery/grpc/pb --go_opt paths=source_relative \
	--go-grpc_out=internal/delivery/grpc/pb --go-grpc_opt paths=source_relative \
	--grpc-gateway_out=internal/delivery/grpc/pb --grpc-gateway_opt paths=source_relative \
	--openapiv2_out=docs/swagger --openapiv2_opt allow_merge=true,merge_file_name=simplebank \
	internal/delivery/grpc/proto/*.proto

rerun_compose:
	docker compose down &
	docker compose up --build

evans:
	evans -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7.2.4-alpine

redis_healthcheck:
	docker exec -it redis redis-cli ping

.PHONY: postgres createdb dropdb migrate-up migrate-down sqlc test server mock dbdocs proto rerun_compose redis evans redis_healthcheck