# Payment Flow

Pre-authorization model (Uber/Bolt style): funds are held before driver search and only charged on trip completion.

## Full Flow

```
Rider                    API Gateway              Trip Service         Payment Service         Driver Service
  |                          |                        |                     |                      |
  |-- POST /trip/start ------>|                        |                     |                      |
  |                          |-- CreateTrip (gRPC) -->|                     |                      |
  |                          |<-- {tripID, amount} ---|                     |                      |
  |                          |                        |                     |                      |
  |                          |-- CreatePaymentIntent (gRPC) -------------->|                      |
  |                          |                        |         Stripe: PaymentIntent              |
  |                          |                        |         capture_method: manual             |
  |                          |<-- {clientSecret} ----------------------------                      |
  |                          |                        |                     |                      |
  |<-- {tripID, clientSecret}|                        |                     |                      |
  |                          |                        |                     |                      |
  | [Stripe Elements form]   |                        |                     |                      |
  | confirmCardPayment()     |                        |                     |                      |
  | → funds AUTHORIZED (held)|                        |                     |                      |
  |                          |                        |                     |                      |
  |-- WS: rider.cmd.payment_confirmed ---------------->|                    |                      |
  |                          |-- GetTripByID (gRPC) ->|                     |                      |
  |                          |-- publish trip.event.created (RabbitMQ) ---------------------------->|
  |<-- WS: trip.event.created|                        |                     |                      |
  |                          |                        |                     |                      |
  |   [Looking for driver]   |                        |                     | driver search...     |
  |                          |                        |                     |                      |
```

### Branch A — Driver found

```
Driver Service           Trip Service             API Gateway              Rider
     |                       |                        |                      |
     |-- driver.event.driver_notified --------------->|                      |
     |                       |-- UpdateTrip (awaiting_driver)                |
     |                       |                        |                      |
     |-- driver.cmd.trip_accept --------------------->|                      |
     |                       |-- UpdateTrip (assigned)|                      |
     |                       |-- publish trip.event.driver_assigned -------->|
     |                       |                        |<-- WS: driver_assigned
     |                       |                        |                      |
     |   [Driver on the way] |                        |                      |
     |                       |                        |                      |
     |-- WS: driver.cmd.trip_complete --------------->|                      |
     |                       |-- CompleteTrip (gRPC)->|                      |
     |                       |   status → "completed" |                      |
     |                       |-- publish payment.cmd.capture_payment ------->|
     |                       |                        |   Stripe: Capture    |
     |                       |                        |   funds CHARGED      |
     |                       |                        |                      |
     |                       |                        |-- WS: trip.event.completed -->|
     |                       |                        |                      |
     |                       |                        |                   [Trip done]
```

### Branch B — No drivers found

```
Driver Service           Trip Service             Payment Service          Rider
     |                       |                        |                      |
     |-- publish trip.event.no_drivers_found -------->|                      |
     |                       |-- UpdateTrip (cancelled)                      |
     |                       |-- publish payment.cmd.cancel_payment -------->|
     |                       |                        |   Stripe: Cancel     |
     |                       |                        |   funds RELEASED     |
     |                       |                        |                      |
     |                  API Gateway <-- WS: trip.event.no_drivers_found ---->|
     |                       |                        |                 [Go back]
```

## Trip Statuses

| Status | Description |
|--------|-------------|
| `pending` | Trip created, waiting for payment confirmation |
| `awaiting_driver` | Driver notified, waiting for response |
| `assigned` | Driver accepted |
| `completed` | Trip finished, payment captured |
| `cancelled` | No drivers found or rider cancelled, payment hold released |

## Services Involved

| Service | Role |
|---------|------|
| **api-gateway** | HTTP handler, WS broker, gRPC orchestrator |
| **trip-service** | Trip lifecycle, gRPC server `:9093` |
| **payment-service** | Stripe integration, gRPC server `:9004`, DB: `payment_db` |
| **driver-service** | Driver matching, notifications |

## RabbitMQ Routing Keys

| Key | Publisher | Consumer |
|-----|-----------|----------|
| `trip.event.created` | api-gateway (WS handler) | driver-service |
| `trip.event.driver_assigned` | trip-service | api-gateway → rider WS |
| `trip.event.no_drivers_found` | driver-service | trip-service, api-gateway → rider WS |
| `payment.cmd.capture_payment` | trip-service | payment-service |
| `payment.cmd.cancel_payment` | trip-service | payment-service |

## Stripe Operations

| Operation | When | Stripe Call |
|-----------|------|-------------|
| `CreatePaymentIntent` | After trip created | `paymentintent.New` with `capture_method: manual` |
| `CapturePayment` | After trip completed | `paymentintent.Capture` |
| `CancelPayment` | No drivers found | `paymentintent.Cancel` |
