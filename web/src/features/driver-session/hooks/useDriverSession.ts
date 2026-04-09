'use client'

import { useMemo, useState } from 'react'
import * as Geohash from 'ngeohash'

import type { Driver } from '@/features/driver/models/types'
import { useDriverEvents } from '@/features/real-time'
import type { Trip } from '@/features/trip'

type Location = { latitude: number; longitude: number }

export function useDriverSession(userID: string, carId: string, location: Location) {
  const [sessionLocation, setSessionLocation] = useState<Location>(location)

  const geohash = useMemo(
    () => Geohash.encode(sessionLocation.latitude, sessionLocation.longitude, 7),
    [sessionLocation.latitude, sessionLocation.longitude],
  )

  const events = useDriverEvents({
    userID,
    carId,
    latitude: sessionLocation.latitude,
    longitude: sessionLocation.longitude,
  })

  const acceptTrip = (trip: Trip, driver: Driver) => {
    events.sendRaw('driver.cmd.trip_accept', {
      tripID: trip.id,
      riderID: trip.userID,
      driver,
    })
    events.reset()
  }

  const declineTrip = (trip: Trip, driver: Driver) => {
    events.sendRaw('driver.cmd.trip_decline', {
      tripID: trip.id,
      riderID: trip.userID,
      driver,
    })
    events.reset()
  }

  const moveLocation = (lat: number, lng: number) => {
    setSessionLocation({ latitude: lat, longitude: lng })
  }

  return {
    driver: events.driver,
    requestedTrip: events.requestedTrip,
    tripStatus: events.tripStatus,
    timeRemaining: events.timeRemaining,
    error: events.error,
    geohash,
    sessionLocation,
    acceptTrip,
    declineTrip,
    moveLocation,
  }
}
