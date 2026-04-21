import type { Driver } from '@/features/driver/models/types'
import type { Trip } from '@/features/trip'

// All WebSocket event types from the server
export const WsEventType = {
  // Driver location broadcast to riders
  DriverLocation: 'driver.cmd.location',
  // Trip lifecycle
  TripCreated: 'trip.event.created',
  TripCompleted: 'trip.event.completed',
  NoDriversFound: 'trip.event.no_drivers_found',
  DriverAssigned: 'trip.event.driver_assigned',
  // Driver-side events
  DriverTripRequest: 'driver.cmd.trip_request',
  DriverTripRequestExpired: 'driver.cmd.trip_request_expired',
  DriverRegister: 'driver.cmd.register',

} as const

export type WsEventType = typeof WsEventType[keyof typeof WsEventType]

// Payload shapes per event type
export type WsEventPayloadMap = {
  [WsEventType.DriverLocation]: Driver[]
  [WsEventType.TripCreated]: Trip
  [WsEventType.TripCompleted]: { tripID: string }
  [WsEventType.NoDriversFound]: undefined
  [WsEventType.DriverAssigned]: Trip
  [WsEventType.DriverTripRequest]: Trip
  [WsEventType.DriverTripRequestExpired]: undefined
  [WsEventType.DriverRegister]: Driver
}

export type WsMessage<T extends WsEventType = WsEventType> = {
  type: T
  data: WsEventPayloadMap[T]
}
