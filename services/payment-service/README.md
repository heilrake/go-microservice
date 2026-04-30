# payment-service

Event-driven сервіс для обробки платежів через Stripe. Слухає події з RabbitMQ та керує lifecycle payment intent'ів.

## Структура

```
services/payment-service/
├── cmd/main.go
├── internal/
│   ├── domain/           # інтерфейси Service, Repository, PaymentProcessor + моделі
│   ├── service/          # бізнес-логіка + integration tests
│   ├── migrations/       # SQL міграції (embed у бінарник)
│   └── infrastructure/
│       ├── repository/   # PostgreSQL через GORM
│       ├── stripe/       # Stripe SDK клієнт
│       ├── events/       # RabbitMQ consumer (trip events)
│       └── grpc/         # gRPC handler
└── pkg/types/            # публічні типи для інших сервісів
```

`service` залежить виключно від `domain` інтерфейсів — жодних прямих залежностей від Stripe чи БД.

## Payment Intent Lifecycle

```
trip.event.created    → CreatePaymentIntent → status: authorized
trip.event.completed  → CapturePayment      → status: captured
trip.event.cancelled  → CancelPayment       → status: cancelled
```

Статуси зберігаються в таблиці `payment_intents` (PostgreSQL).

## Integration Tests

Файл: `internal/service/service_integration_test.go`

Тести запускаються проти реального PostgreSQL у Docker через [testcontainers-go](https://golang.testcontainers.org/). Stripe замінений `mockStripe` — без реального API.

### Як це працює

1. `TestMain` піднімає PostgreSQL контейнер (`postgres:16`)
2. Застосовує міграції через `sharedBootstrap.RunMigrator`
3. Підключається через GORM і передає `*gorm.DB` у тести
4. Кожен тест викликає `cleanDB(t)` для ізоляції між собою

### Що покрито

| Тест | Сценарій |
|------|----------|
| `TestCreatePaymentIntent` | Створення intent, статус `authorized`, збереження в БД |
| `TestCapturePayment` | Захоплення платежу → статус `captured` |
| `TestCancelPayment` | Скасування → статус `cancelled` |
| `TestCapturePayment_TripNotFound` | Помилка якщо trip не існує |
| `TestCancelPayment_TripNotFound` | Помилка якщо trip не існує |
| `TestCreatePaymentIntent_DuplicateTripID` | UNIQUE constraint: другий create → помилка |
| `TestCapturePayment_StripeError_StatusUnchanged` | Stripe впав → статус лишається `authorized` |
| `TestCapturePayment_ConcurrentCalls` | 5 горутин одночасно → фінальний статус `captured` |
| `TestCreatePaymentIntent_ConcurrentDuplicates` | 5 горутин одночасно → рівно 1 успіх |

### Запуск

```bash
# з кореня репозиторію
go test ./services/payment-service/internal/service/... -v -timeout 120s
```

> Потрібен запущений Docker.
