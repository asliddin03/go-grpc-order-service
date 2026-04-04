.PHONY: up down logs test test-unit test-integration build proto proto-order proto-auth proto-inventory

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

test: test-unit

test-unit:
	cd order-service && go test ./internal/...
	cd auth-service && go test ./...
	cd inventory-service && go test ./...

test-integration:
	cd order-service && POSTGRES_DSN="$(POSTGRES_DSN)" go test ./tests/integration/...

build:
	cd order-service && go build ./...
	cd auth-service && go build ./...
	cd inventory-service && go build ./...

proto: proto-order proto-auth proto-inventory

proto-order:
	cd order-service && protoc --proto_path=proto --go_out=. --go-grpc_out=. proto/order.proto

proto-auth:
	cd auth-service && protoc --proto_path=proto --go_out=. --go-grpc_out=. proto/auth.proto

proto-inventory:
	cd inventory-service && protoc --proto_path=proto --go_out=. --go-grpc_out=. proto/inventory.proto