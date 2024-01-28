postgres:
	docker run --name postgres-container -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234 -d postgres:alpine3.18

createdb:
	docker exec -it postgres-container createdb --username=root --owner=root messenger

dropdb:
	docker exec -it postgres-container dropdb messenger

createmigrationfile:
	migrate create -ext sql -dir db/migration -seq init_schema

migrateup:
	 migrate -path ./db/migration -database "postgres://root:1234@localhost:5432/messenger?sslmode=disable" -verbose up

migratedown:
	 migrate -path ./db/migration -database "postgres://root:1234@localhost:5432/messenger?sslmode=disable" -verbose down

test:
	go test -v -cover ./...

.PHONY: postgres, createdb, dropdb, migratedown, migrateup