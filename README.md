# go-grpc-order-service

Небольшой backend-проект на Go с микросервисной архитектурой.  
Проект демонстрирует взаимодействие сервисов через gRPC, хранение данных в PostgreSQL, слоистую архитектуру, unit tests и integration tests.

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

## Сценарий взаимодействия

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