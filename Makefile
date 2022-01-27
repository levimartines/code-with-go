postgres:
	docker-compose up db

migrateup:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5433/postgres?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5433/postgres?sslmode=disable" -verbose down

sqlc:
	sqlc generate