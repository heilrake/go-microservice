import type { Driver } from "@/features/driver/models/types";
import type { Coordinate, Route, RouteFare, Trip } from "@/features/trip";


// These are the endpoints the API Gateway must have for the frontend to work correctly
export const BackendEndpoints = {
  PREVIEW_TRIP: "/trip/preview",
  START_TRIP: "/trip/start",
  WS_DRIVERS: "/drivers",
  WS_RIDERS: "/riders",
} as const;

export const TripEvents = {
  NoDriversFound: "trip.event.no_drivers_found",
  DriverAssigned: "trip.event.driver_assigned",
  Completed: "trip.event.completed",
  Cancelled: "trip.event.cancelled",
  Created: "trip.event.created",
  DriverLocation: "driver.cmd.location",
  DriverTripRequest: "driver.cmd.trip_request",
  DriverTripAccept: "driver.cmd.trip_accept",
  DriverTripDecline: "driver.cmd.trip_decline",
  DriverRegister: "driver.cmd.register",
  PaymentSessionCreated: "payment.event.session_created",
} as const;

export type TripEventType = typeof TripEvents[keyof typeof TripEvents];

// Messages sent from the server to the client via the websocket
export type ServerWsMessage =
  | PaymentSessionCreatedRequest
  | DriverAssignedRequest
  | DriverLocationRequest
  | DriverTripRequest
  | DriverRegisterRequest
  | TripCreatedRequest
  | NoDriversFoundRequest;

// Messages sent from the client to the server via the websocket
export type ClientWsMessage = DriverResponseToTripResponse

type TripCreatedRequest = {
  type: typeof TripEvents.Created;
  data: Trip;
}

type NoDriversFoundRequest = {
  type: typeof TripEvents.NoDriversFound;
}

type DriverRegisterRequest = {
  type: typeof TripEvents.DriverRegister;
  data: Driver;
}
type DriverTripRequest = {
  type: typeof TripEvents.DriverTripRequest;
  data: Trip;
}

export type PaymentEventSessionCreatedData = {
  tripID: string;
  sessionID: string;
  amount: number;
  currency: string;
};

type PaymentSessionCreatedRequest = {
  type: typeof TripEvents.PaymentSessionCreated;
  data: PaymentEventSessionCreatedData;
};

type DriverAssignedRequest = {
  type: typeof TripEvents.DriverAssigned;
  data: Trip;
}

type DriverLocationRequest = {
  type: typeof TripEvents.DriverLocation;
  data: Driver[];
}

type DriverResponseToTripResponse = {
  type: typeof TripEvents.DriverTripAccept | typeof TripEvents.DriverTripDecline;
  data: {
    tripID: string;
    riderID: string;
    driver: Driver;
  };
}

export type HTTPTripPreviewResponse = {
  route: Route;
  rideFares: RouteFare[];
}

export type HTTPTripStartRequestPayload = {
  rideFareID: string;
  userID: string;
}

export type HTTPTripPreviewRequestPayload = {
  userID: string;
  pickup: Coordinate;
  destination: Coordinate;
}

export function isValidTripEvent(event: string): event is TripEventType {
  return Object.values(TripEvents).includes(event as TripEventType);
}

export function isValidWsMessage(
  message: { type: string }
): message is ServerWsMessage {
  return isValidTripEvent(message.type);
}
