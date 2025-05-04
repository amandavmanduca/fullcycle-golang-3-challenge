include .env

setup:
	docker compose up -d --build
	$(MAKE) check-mysql
	@if [ "$(SKIP_MIGRATIONS)" != "true" ]; then \
		$(MAKE) migrate-up; \
	fi

check-mysql:
	@echo "Aguardando MySQL estar pronto na porta $(MYSQL_PORT)..."
	@until nc -z $(MYSQL_HOST) $(MYSQL_PORT) && echo "MySQL pronto"; do \
		sleep 2; \
		echo "Esperando MySQL..."; \
	done
	@echo "MySQL está pronto."

check-rabbitmq:
	@echo "Aguardando RabbitMQ responder na interface de gerenciamento..."
	@until curl -s -u $(RABBITMQ_USER):$(RABBITMQ_PASS) http://localhost:$(RABBITMQ_MANAGEMENT_PORT)/api/healthchecks/node | grep -q '"status":"ok"'; do \
		sleep 1; \
	done

check:
	@if ! docker compose ps | grep -q 'Up'; then \
		echo "Containers não estão rodando. Executando setup..."; \
		$(MAKE) setup; \
	fi
	$(MAKE) check-mysql
	$(MAKE) check-rabbitmq


run: check
	cd cmd/ordersystem && go run main.go

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

setup-down:
	docker-compose down