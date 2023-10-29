
storage/library/queries/db.go: storage/library/sqlc.yaml storage/library/queries/*.sql storage/library/migrations/*.sql
	@sqlc generate -f ./storage/library/sqlc.yaml

.PHONY: build
build:  storage/library/queries/db.go
	@go build .

run: build
	./main
