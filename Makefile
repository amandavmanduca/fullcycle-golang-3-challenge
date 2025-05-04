include .env

setup:
	docker compose up -d --build

check:
	@if ! docker compose ps | grep -q 'Up'; then \
		echo "Containers não estão rodando. Executando setup..."; \
		$(MAKE) setup; \
	fi

	@echo "Aguardando MySQL estar pronto na porta $(MYSQL_PORT)..."
	@until nc -z localhost $(MYSQL_PORT); do \
		sleep 1; \
	done
	@echo "MySQL está pronto."

	@echo "Aguardando RabbitMQ responder na interface de gerenciamento..."
	@until curl -s -u $(RABBITMQ_USER):$(RABBITMQ_PASS) http://localhost:$(RABBITMQ_MANAGEMENT_PORT)/api/healthchecks/node | grep -q '"status":"ok"'; do \
		sleep 1; \
	done
	@echo "RabbitMQ está pronto."

run: check
	cd cmd/ordersystem && go run main.go

gen-graphql:
	go run github.com/99designs/gqlgen generate

gen-proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && protoc --go_out=. --go-grpc_out=. internal/infra/grpc/protofiles/*

evans:
	go install github.com/ktr0731/evans@latest && evans -r repl

# go install github.com/ktr0731/evans@latest