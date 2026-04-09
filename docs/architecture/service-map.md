# Service Communication Map

> Update this diagram whenever a new service, gRPC method, or RabbitMQ event is added.

```mermaid
flowchart TD
    Client["Client\n(Browser / Mobile)"]
    Stripe["Stripe"]

    subgraph GW ["API Gateway :8081"]
        REST["REST Handlers"]
        WS["WebSocket Hub"]
        SWAGGER["Swagger UI /swagger/"]
    end

    subgraph AuthSvc ["Auth Service :8082"]
        OAuth["OAuth / JWT"]
    end

    subgraph gRPC ["gRPC Services (internal)"]
        UserSvc["User Service"]
        DriverSvc["Driver Service"]
        TripSvc["Trip Service"]
        PaySvc["Payment Service"]
    end

    MQ["RabbitMQ :5672\n(UI :15672)"]

    %% Client connections
    Client -->|"HTTP REST"| REST
    Client <-->|"WebSocket\n/ws/drivers  /ws/riders"| WS
    Stripe -->|"POST /webhook/stripe"| REST

    %% Auth proxy
    REST -->|"HTTP proxy\nPOST /auth/oauth\nPOST /dev/login"| OAuth
    OAuth -->|"gRPC: LoginUser\nGetOrCreateUserByOAuth"| UserSvc

    %% gRPC from gateway
    REST -->|"gRPC: CreateUser"| UserSvc
    REST -->|"gRPC: CreateDriver\nGetDriver\nCreateCar\nListCars"| DriverSvc
    REST -->|"gRPC: PreviewTrip\nCreateTrip"| TripSvc

    %% RabbitMQ publish
    REST -->|"publish: payment.success"| MQ

    %% WebSocket hub: routes MQ messages to connected clients by ownerID
    WS -->|"consume: driver.cmd.trip_request\nnotify.driver.no_drivers\nnotify.driver.assignment\nnotify.payment.session"| MQ

    %% Driver WebSocket → RabbitMQ (accept / decline)
    WS -->|"publish: driver.cmd.trip_accept\ndriver.cmd.trip_decline"| MQ

    %% Internal service events
    TripSvc -->|"publish: trip.event.created"| MQ
    REST -->|"WS push: trip.event.created → rider"| WS
    PaySvc  -->|"publish: payment.event.session_created"| MQ
    MQ      -->|"consume"| PaySvc

    %% Driver-service: find & notify flow
    MQ -->|"find_available_drivers\n(trip.event.created\ntrip.event.driver_not_interested)"| DriverSvc
    DriverSvc -->|"publish: driver.cmd.trip_request"| MQ
    DriverSvc -->|"publish: driver.event.driver_notified"| MQ

    %% Trip-service: driver response flow
    MQ -->|"driver_trip_response\n(driver.cmd.trip_accept\ndriver.cmd.trip_decline)"| TripSvc
    MQ -->|"driver_notified\n(driver.event.driver_notified)"| TripSvc
    MQ -->|"trip_search_failed\n(trip.event.no_drivers_found)"| TripSvc
    TripSvc -->|"publish: trip.event.driver_assigned\npayment.cmd.create_session\ntrip.event.driver_not_interested"| MQ

    %% Driver-service: ack queue (cancel 15s timer)
    MQ -->|"driver_trip_ack\n(driver.cmd.trip_accept\ndriver.cmd.trip_decline)"| DriverSvc

    %% API Gateway routes expired notification to driver via WebSocket
    DriverSvc -->|"publish: driver.cmd.trip_request_expired"| MQ
    MQ -->|"driver_trip_request_expired"| GW
```

## Driver Accept/Decline Flow (15s timeout)

```mermaid
sequenceDiagram
    participant Rider
    participant TripSvc as Trip Service
    participant MQ as RabbitMQ
    participant DriverSvc as Driver Service
    participant GW as API Gateway
    participant Driver

    Rider->>TripSvc: CreateTrip (via GW REST)
    GW->>Rider: WS push: trip.event.created → UI shows spinner
    TripSvc->>MQ: trip.event.created [status: pending]
    MQ->>DriverSvc: find_available_drivers
    DriverSvc->>MQ: driver.cmd.trip_request (ownerID=driverUserID)
    DriverSvc->>MQ: driver.event.driver_notified {tripID}
    MQ->>GW: driver.cmd.trip_request → WS
    GW->>Driver: WebSocket: trip request
    MQ->>TripSvc: driver_notified → status: awaiting_driver

    alt Driver accepts within 15s
        Driver->>GW: WS: driver.cmd.trip_accept
        GW->>MQ: driver.cmd.trip_accept
        MQ->>DriverSvc: driver_trip_ack → cancel timer
        MQ->>TripSvc: driver_trip_response → status: assigned
        TripSvc->>MQ: trip.event.driver_assigned + payment.cmd.create_session
    else No drivers available / max retries exceeded
        DriverSvc->>MQ: trip.event.no_drivers_found {tripID, riderID}
        MQ->>GW: notify_driver_no_drivers → WS (Rider UI: "No drivers found")
        MQ->>TripSvc: trip_search_failed → status: cancelled
    else Driver declines within 15s
        Driver->>GW: WS: driver.cmd.trip_decline
        GW->>MQ: driver.cmd.trip_decline
        MQ->>DriverSvc: driver_trip_ack → cancel timer
        MQ->>TripSvc: driver_trip_response → re-publish trip.event.driver_not_interested
        MQ->>DriverSvc: find_available_drivers → next driver
    else 15s timeout (no response)
        DriverSvc->>MQ: driver.cmd.trip_request_expired (ownerID=driverUserID)
        DriverSvc->>MQ: trip.event.driver_not_interested
        MQ->>GW: driver_trip_request_expired → WS
        GW->>Driver: WebSocket: request expired → UI resets
        MQ->>DriverSvc: find_available_drivers → next driver
    end
```

## Ports

| Service        | Port  | Protocol     |
|----------------|-------|--------------|
| API Gateway    | 8081  | HTTP / WS    |
| Swagger UI     | 8081  | HTTP (`/swagger/`) |
| Auth Service   | 8082  | HTTP (internal) |
| Proto Docs     | 8090  | HTTP (Tilt)  |
| RabbitMQ       | 5672  | AMQP         |
| RabbitMQ UI    | 15672 | HTTP         |
| Jaeger UI      | 16686 | HTTP         |
| Trip Postgres  | 30432 | TCP          |
| Driver Postgres| 30433 | TCP          |
| User Postgres  | 30434 | TCP          |

## RabbitMQ Queues

| Queue                          | Routing key(s)                                          | Producer        | Consumer                      |
|--------------------------------|---------------------------------------------------------|-----------------|-------------------------------|
| `find_available_drivers`       | `trip.event.created`, `trip.event.driver_not_interested`| Trip Svc / Driver Svc | Driver Service          |
| `driver_cmd_trip_request`      | `driver.cmd.trip_request`                               | Driver Service  | API GW → Driver WS            |
| `driver_trip_response`         | `driver.cmd.trip_accept`, `driver.cmd.trip_decline`     | API Gateway     | Trip Service                  |
| `driver_trip_ack`              | `driver.cmd.trip_accept`, `driver.cmd.trip_decline`     | API Gateway     | Driver Service (cancel timer) |
| `driver_notified`              | `driver.event.driver_notified`                          | Driver Service  | Trip Service (status update)  |
| `driver_trip_request_expired`  | `driver.cmd.trip_request_expired`                       | Driver Service  | API GW → Driver WS → reset UI |
| `notify_driver_no_drivers`     | `trip.event.no_drivers_found`                           | Driver Service  | API GW → Rider WS             |
| `trip_search_failed`           | `trip.event.no_drivers_found`                           | Driver Service  | Trip Service (status: cancelled) |
| `notify_driver_assignment`     | `trip.event.driver_assigned`                            | Trip Service    | API GW → Rider WS             |
| `notify_payment_session`       | `payment.event.session_created`                         | Payment Service | API GW → Rider WS             |
| `payment_trip_response`        | `payment.cmd.create_session`                            | Trip Service    | Payment Service               |

## Trip Status Flow

| Status             | Transition                                          |
|--------------------|-----------------------------------------------------|
| `pending`          | Trip created                                        |
| `awaiting_driver`  | Driver notified (`driver.event.driver_notified`)    |
| `awaiting_driver`  | Next driver notified (after decline / timeout)      |
| `assigned`         | Driver accepted (`driver.cmd.trip_accept`)          |
| `cancelled`        | No drivers found (`trip.event.no_drivers_found`)    |

## gRPC Methods

### UserService
| Method                    | Request fields                                      |
|---------------------------|-----------------------------------------------------|
| `CreateUser`              | username, email, password, role, profile_picture    |
| `UpdateUser`              | user_id, username?, email?, profile_picture?        |
| `LoginUser`               | email, password, role                               |
| `GetUser`                 | user_id                                             |
| `GetOrCreateUserByOAuth`  | email, username, profile_picture, role              |

### DriverService
| Method            | Request fields                                    |
|-------------------|---------------------------------------------------|
| `CreateDriver`    | user_id, name, profile_picture                    |
| `GetDriver`       | user_id                                           |
| `CreateCar`       | user_id, car_plate, package_slug                  |
| `ListCars`        | user_id                                           |
| `RegisterDriver`  | driverID (user_id), car_id, latitude, longitude   |
| `UnRegisterDriver`| driverID (user_id), car_id, latitude, longitude   |

### TripService
| Method        | Request fields                              |
|---------------|---------------------------------------------|
| `PreviewTrip` | userID, startLocation{lat,lon}, endLocation{lat,lon} |
| `CreateTrip`  | rideFareID, userID                          |
