# TODO & Code Review

## Баги / проблеми що треба виправити

### 1. `trip-service` — debug код і непрацюючий error handling
**Файл:** `services/trip-service/internal/infrastructure/http/http.go:33`
```go
// ЗАРАЗ (зламано):
t, err := s.Service.GetRoute(ctx, &requestBody.Pickup, &requestBody.Destination)
if err != nil {
    fmt.Println("hello dude")  // <- debug залишок, помилка ігнорується!
}
writeJson(w, http.StatusOK, t) // <- завжди 200, навіть якщо err != nil
```
Треба повернути `http.StatusInternalServerError` і прибрати `fmt.Println`.

---

### 2. `trip-service` — HTTP клієнт без таймауту
**Файл:** `services/trip-service/internal/service/service.go:73`
```go
resp, err := http.Get(url) // <- дефолтний клієнт, немає таймауту
```
OSRM може висіти вічно. Треба `http.Client{Timeout: 5 * time.Second}`.

---

### 3. `driver-service` — infrastructure тип протікає через шари
**Файл:** `services/driver-service/internal/service/service.go`

Методи `CreateDriver`, `GetDriver`, `CreateCar`, `ListCars`, `RegisterDriver`, `FindAvailableDrivers` повертають `*infrastructure.DriverModel` / `*infrastructure.CarModel` напряму з service layer. Це порушує Clean Architecture — service не повинен знати про DB-моделі. Треба domain-моделі у `internal/domain/`.

---

### 4. `trip-service/repository` — застарілий порівняння помилки
**Файл:** `services/trip-service/internal/infrastructure/repository/postgres.go:77`
```go
if err == gorm.ErrRecordNotFound { // <- треба errors.Is()
```
Те саме в `GetTripByID`. GORM v2 рекомендує `errors.Is(err, gorm.ErrRecordNotFound)`.

---

### 5. `driver-service` — `maxSearchAttempts = 1`
**Файл:** `services/driver-service/internal/events/trip_consumer.go:89`
```go
const maxSearchAttempts = 1
```
Одна спроба — якщо водій не відповів, поїздку одразу скасовано. Скоріше за все треба 3-5. Або зробити конфігурованим через env.

---

### 6. `payment-service` — routing key не перевіряється в `ListenCapture` / `ListenCancel`
**Файл:** `services/payment-service/internal/infrastructure/events/trip_consumer.go`

`Listen()` перевіряє `msg.RoutingKey` через switch. `ListenCapture()` і `ListenCancel()` — ні, одразу обробляють будь-яке повідомлення з черги. Якщо в чергу потрапить чужий routing key — буде тихий баг.

---

### 7. `CancelTrip` — приймає `userID` замість `tripID`
**Файл:** `services/trip-service/internal/infrastructure/repository/postgres.go:39`

Функція `CancelTrip(ctx, userID)` скасовує всі активні поїздки юзера, а не конкретну поїздку. Це може бути навмисне, але назва та сигнатура вводять в оману (`tripID string` в інтерфейсі vs реальна логіка).

---

## Технічний борг

### 8. Немає тестів для trip-service, driver-service, user-service
Є тільки для payment-service. Найважливіші для покриття:
- **trip-service:** `EstimatePackagesPriceWithRoute`, `GetAndValidateFare`, OSRM інтеграція
- **driver-service:** `FindAvailableDrivers` (geohash + JOIN логіка), `RegisterDriver` / `UnregisterDriver`
- **user-service:** `GetOrCreateUserByOAuth` (find-or-create race condition)

---

### 9. Немає DLQ consumer
Dead letter queue для payment capture є (налаштована в RabbitMQ), але consumer який її читає — відсутній. Повідомлення що failed після retry просто накопичуються, ніхто не сповіщає, не алертить, не ретраїть.

---

### 10. Structured logging відсутній
Скрізь використовується `log.Printf`. Немає рівнів (debug/info/warn/error), немає контекстних полів (trip_id, user_id тощо). Варто підключити `slog` (стандартний з Go 1.21) або `zap`.

---

### 11. OpenTelemetry тільки в driver-service
`driver-service` має повноцінний tracing (spans, attributes, links між goroutine). В `trip-service`, `payment-service`, `user-service` — нічого. Розподілені traces обриваються на межі сервісів.

---

### 12. `bootstrap.InitGorm` — `log.Fatalf` замість повернення помилки
**Файл:** `shared/bootstrap/postgres.go`

`log.Fatalf` в бібліотечному коді — погана практика. Функція `InitGorm` не дає caller'у шансу обробити помилку. Треба повертати `(*gorm.DB, error)`.

---

### 13. `getBaseFares()` — hardcoded тарифи
**Файл:** `services/trip-service/internal/service/service.go:152`

Базові ціни (`suv: 200`, `sedan: 350`, і т.д.) захардкоджені в коді. Зміна тарифу = деплой. Варто винести в конфіг або БД.

---

### 14. `user-service` — мертвий код `LoginUser` / `CreateUser` з password
Пам'ять каже що email/password автентифікацію видалено, але в `user-service/service.go` є `LoginUser` з bcrypt та `CreateUser` з хешуванням паролю. Якщо вони більше не викликаються через API — це мертвий код.

---

## Що можна додати / покращити

### A. Integration tests — driver-service
Найцікавіше для тестування:
- `RegisterDriver` → `FindAvailableDrivers` → `UnregisterDriver` (повний lifecycle)
- Конкурентна реєстрація кількох водіїв
- `FindAvailableDrivers` з різними `packageSlug`

### B. Integration tests — trip-service
- `CreateTrip` + `UpdateTrip` → `CompleteTrip`
- `CancelTrip` при різних статусах
- `GetAndValidateFare` — чужий userID

### C. gRPC health checks
Стандартний `google.golang.org/grpc/health` протокол — дозволяє Kubernetes liveness/readiness probам нормально перевіряти сервіси.

### D. Геопошук водіїв
Зараз `FindAvailableDrivers` повертає всіх доступних водіїв з потрібним `packageSlug` по всьому місту. Geohash є в БД але не використовується для фільтрації. Можна шукати тільки в радіусі N км від pickup точки.

### E. Trip history endpoint
Немає API щоб отримати минулі поїздки rider'а. `GET /trips?userID=` або `GET /me/trips`.

### F. Webhook від Stripe (або polling)
Зараз payment capture викликається через RabbitMQ event від trip-service. Але якщо capture фейлиться — немає механізму retry окрім DLQ. Stripe webhooks дозволяють отримати підтвердження успіху/фейлу напряму.

---

## Пріоритети (суб'єктивно)

| # | Завдання | Важливість |
|---|----------|------------|
| 1 | Виправити `fmt.Println("hello dude")` + error handling в http.go | критично |
| 2 | HTTP таймаут для OSRM | критично |
| 3 | `errors.Is` замість `==` для GORM | середнє |
| 4 | Domain моделі в driver-service (рефакторинг) | середнє |
| 5 | DLQ consumer | середнє |
| 6 | Integration tests для driver-service | середнє |
| 7 | `maxSearchAttempts` конфігурабельний | низьке |
| 8 | Structured logging (slog) | низьке |
| 9 | OpenTelemetry у всіх сервісах | низьке |
| 10 | Геопошук водіїв | фіча |
| 11 | Trip history | фіча |
