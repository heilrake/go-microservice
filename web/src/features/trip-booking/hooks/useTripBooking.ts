'use client'

import { useRef, useState } from 'react'

import type { CarPackageSlugType } from '@/features/packages'
import type { RouteFare, TripPreview } from '@/features/trip'

import { tripApi } from '@/shared/api'

type Location = { latitude: number; longitude: number }

export function useTripBooking(userID: string, location: Location) {
  const [trip, setTrip] = useState<TripPreview | null>(null)
  const [destination, setDestination] = useState<[number, number] | null>(null)
  const [selectedCarType, setSelectedCarType] = useState<CarPackageSlugType | null>(null)
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  const requestPreview = async (destinationCoords: [number, number]) => {
    const data = await tripApi.preview({
      userID,
      pickup: { latitude: location.latitude, longitude: location.longitude },
      destination: { latitude: destinationCoords[0], longitude: destinationCoords[1] },
    })

    const parsedRoute = data.route.geometry[0].coordinates
      .map((coord: { longitude: number; latitude: number }) => [coord.longitude, coord.latitude] as [number, number])

    setTrip({
      tripID: '',
      route: parsedRoute,
      rideFares: data.rideFares,
      distance: data.route.distance,
      duration: data.route.duration,
    })
  }

  const handleMapClick = (lat: number, lng: number) => {
    if (trip?.tripID) return

    if (debounceRef.current) clearTimeout(debounceRef.current)

    debounceRef.current = setTimeout(async () => {
      setDestination([lat, lng])
      await requestPreview([lat, lng])
    }, 500)
  }

  const startTrip = async (fare: RouteFare) => {
    if (!fare.id) throw new Error('No Fare ID in the payload')

    setSelectedCarType(fare.packageSlug)

    const data = await tripApi.start({
      rideFareID: fare.id,
      userID: fare.userID || userID,
    })

    setTrip(prev => prev ? { ...prev, tripID: data.tripID } : prev)

    return data
  }

  const reset = () => {
    setTrip(null)
    setDestination(null)
    setSelectedCarType(null)
  }

  return { trip, destination, selectedCarType, handleMapClick, startTrip, reset }
}
