postgres :
	docker run --name postgres14.2 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=0237 -d postgres:14.2-alpine
createdb: 
	docker exec -it postgres14.2 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres14.2 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:0237@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:0237@localhost:5432/simple_bank?sslmode=disable" -verbose down 

sqlc:
	sqlc generate

test: 
	go clean -testcache && go test -v -cover ./...

server: 
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server