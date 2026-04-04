# go-grpc-order-service

Небольшой backend-проект на Go с микросервисной архитектурой.  
Проект демонстрирует взаимодействие сервисов через gRPC, хранение данных в PostgreSQL, слоистую архитектуру, unit tests и integration tests.

## Статус проекта

**MVP завершен**

Сейчас в проекте реализованы:
- `order-service` как основной сервис заказов
- `auth-service` как mock gRPC сервис проверки пользователя
- `inventory-service` как mock gRPC сервис товаров
- создание и чтение заказов
- межсервисное взаимодействие через gRPC
- PostgreSQL repository
- unit и integration tests

---

## Обзор проекта

Репозиторий содержит три сервиса:

- **order-service** — основной сервис, который отвечает за создание и получение заказов
- **auth-service** — mock gRPC сервис для проверки пользователя
- **inventory-service** — mock gRPC сервис для получения информации о товарах

Основная цель проекта — показать, как Go-микросервис может:

- предоставлять gRPC API
- взаимодействовать с другими сервисами через gRPC
- хранить данные в PostgreSQL
- использовать чистую слоистую архитектуру
- покрываться unit и integration тестами

---

## Архитектура

### Сервисы

#### order-service
Отвечает за:
- создание заказов
- получение заказа по id
- получение списка заказов пользователя
- проверку пользователя через `auth-service`
- получение данных о товарах через `inventory-service`
- сохранение заказов в PostgreSQL

#### auth-service
Отвечает за:
- проверку существования пользователя
- возврат mock-данных пользователя

#### inventory-service
Отвечает за:
- возврат mock-данных о товарах
- предоставление цены и доступности товара

---

## Схема взаимодействия

```text
Client
  |
  v
order-service
  |----> auth-service
  |----> inventory-service
  |
  v
PostgreSQL
```

---

## Сценарий работы

### CreateOrder
1. Клиент отправляет `CreateOrder` в `order-service`
2. `order-service` вызывает `auth-service` для проверки пользователя
3. `order-service` вызывает `inventory-service` для получения цены и доступности товаров
4. `order-service` считает итоговую сумму заказа
5. `order-service` сохраняет заказ и позиции заказа в PostgreSQL
6. `order-service` возвращает созданный заказ

### Read flow
- `GetOrder` читает один заказ из PostgreSQL
- `ListOrders` читает все заказы пользователя из PostgreSQL

---

## Стек технологий

- **Go**
- **gRPC**
- **Protocol Buffers**
- **PostgreSQL**
- **pgx / pgxpool**
- **Docker Compose**
- **Testify**
- **Makefile**

---

## Структура проекта

```text
go-grpc-order-service/
├── .env.example
├── Makefile
├── docker-compose.yml
├── go.work
├── auth-service/
├── inventory-service/
└── order-service/
```

### Структура `order-service`

```text
order-service/
├── cmd/order/main.go
├── gen/order/v1/
├── internal/
│   ├── app/
│   ├── client/
│   ├── config/
│   ├── domain/
│   ├── handler/grpc/
│   ├── interceptor/
│   ├── repository/postgres/
│   ├── service/
│   └── storage/
├── migrations/
├── proto/
└── tests/integration/
```

---

## Что реализовано

### order-service
- `CreateOrder`
- `GetOrder`
- `ListOrders`
- PostgreSQL repository
- gRPC handlers
- gRPC server
- logging interceptor
- recovery interceptor
- внешние gRPC клиенты для auth и inventory

### auth-service
- `ValidateUser`

### inventory-service
- `GetProducts`

---

## Переменные окружения

Пример значений:

```env
POSTGRES_DSN=postgres://postgres:postgres@localhost:5432/order_service?sslmode=disable
GRPC_PORT=50051
AUTH_SERVICE_ADDRESS=localhost:50052
INVENTORY_SERVICE_ADDRESS=localhost:50053
```

---

## Запуск проекта

### 1. Поднять PostgreSQL

Из корня репозитория:

```bash
make up
```

### 2. Установить переменные окружения

```bash
export POSTGRES_DSN="postgres://postgres:postgres@localhost:5432/order_service?sslmode=disable"
export GRPC_PORT="50051"
export AUTH_SERVICE_ADDRESS="localhost:50052"
export INVENTORY_SERVICE_ADDRESS="localhost:50053"
```

### 3. Запустить auth-service

В отдельном терминале:

```bash
cd auth-service
go run ./cmd/auth
```

### 4. Запустить inventory-service

В отдельном терминале:

```bash
cd inventory-service
go run ./cmd/inventory
```

### 5. Запустить order-service

В отдельном терминале:

```bash
cd order-service
go run ./cmd/order
```

---

## Генерация protobuf-кода

Из корня репозитория:

```bash
make proto
```

Либо отдельно по сервисам.

### order-service

```bash
cd order-service
protoc --proto_path=proto --go_out=. --go-grpc_out=. proto/order.proto
```

### auth-service

```bash
cd auth-service
protoc --proto_path=proto --go_out=. --go-grpc_out=. proto/auth.proto
```

### inventory-service

```bash
cd inventory-service
protoc --proto_path=proto --go_out=. --go-grpc_out=. proto/inventory.proto
```

---

## Запуск тестов

### Запустить все тесты

Из корня репозитория:

```bash
make test
```

### Запустить только тесты `order-service`

```bash
cd order-service
go test ./...
```

### Запустить только integration tests

```bash
cd order-service
go test ./tests/integration/...
```

---

## Примеры gRPC-запросов

### CreateOrder

Запускать из директории `order-service`:

```bash
grpcurl \
  -plaintext \
  -import-path ./proto \
  -proto order.proto \
  -d '{
    "user_id": 42,
    "items": [
      { "product_id": 1, "quantity": 2 },
      { "product_id": 2, "quantity": 3 }
    ]
  }' \
  localhost:50051 \
  order.v1.OrderService/CreateOrder
```

### GetOrder

```bash
grpcurl \
  -plaintext \
  -import-path ./proto \
  -proto order.proto \
  -d '{"order_id": 1}' \
  localhost:50051 \
  order.v1.OrderService/GetOrder
```

### ListOrders

```bash
grpcurl \
  -plaintext \
  -import-path ./proto \
  -proto order.proto \
  -d '{"user_id": 42}' \
  localhost:50051 \
  order.v1.OrderService/ListOrders
```

### ValidateUser

Запускать из директории `auth-service`:

```bash
grpcurl \
  -plaintext \
  -import-path ./proto \
  -proto auth.proto \
  -d '{"user_id": 42}' \
  localhost:50052 \
  auth.v1.AuthService/ValidateUser
```

### GetProducts

Запускать из директории `inventory-service`:

```bash
grpcurl \
  -plaintext \
  -import-path ./proto \
  -proto inventory.proto \
  -d '{"product_ids": [1, 2, 20]}' \
  localhost:50053 \
  inventory.v1.InventoryService/GetProducts
```

---

## Стратегия тестирования

### Unit tests
Unit tests покрывают слой `service` в `order-service`:
- успешное создание заказа
- ошибки валидации
- обработку ошибок зависимостей
- получение заказа
- получение списка заказов

### Integration tests
Integration tests покрывают PostgreSQL repository:
- создание заказа
- получение заказа по id
- получение списка заказов по user id

---

## Что демонстрирует этот проект

- проектирование gRPC API через protobuf
- реализацию слоистой архитектуры на Go
- использование интерфейсов для инверсии зависимостей
- написание unit tests с fake-реализациями
- написание integration tests для PostgreSQL
- взаимодействие между сервисами через gRPC
- организацию небольшого multi-service проекта на Go
