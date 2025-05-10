include .env

run:
	docker compose up -d --build

gen-graphql:
	go run github.com/99designs/gqlgen generate

gen-proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && protoc --go_out=. --go-grpc_out=. internal/infra/grpc/protofiles/*

evans:
	go install github.com/ktr0731/evans@latest && evans -r repl

# make new-migration NAME=desired_name
new-migration:
	go install github.com/pressly/goose/v3/cmd/goose@latest && goose --dir internal/infra/database/migrations create $(NAME) sql

migrate-up:
	@echo "Executing migrations up"
	goose --dir internal/infra/database/migrations mysql "$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)" up

migrate-down:
	@echo "Executing migrations down"
	goose --dir internal/infra/database/migrations mysql "$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)" down

teardown:
	docker-compose down