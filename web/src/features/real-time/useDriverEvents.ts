'use client'

import { useEffect, useRef, useState } from 'react'

import type { Driver } from '@/features/driver/models/types'
import type { Trip } from '@/features/trip'

import { WEBSOCKET_URL } from '@/shared/libs/constants'
import { BackendEndpoints } from '@/shared/libs/contracts'
import { WsClient, WsEventType } from '@/shared/ws'

type DriverEventsState = {
  driver: Driver | null
  requestedTrip: Trip | null
  assignedTrip: Trip | null
  tripStatus: WsEventType | null
  timeRemaining: number | null
  error: string | null
}

type ConnectParams = {
  userID: string
  carId: string
  latitude: number
  longitude: number
}

export function useDriverEvents({ userID, carId, latitude, longitude }: ConnectParams) {
  const clientRef = useRef<WsClient | null>(null)
  const countdownRef = useRef<ReturnType<typeof setInterval> | null>(null)
  const [state, setState] = useState<DriverEventsState>({
    driver: null,
    requestedTrip: null,
    assignedTrip: null,
    tripStatus: null,
    timeRemaining: null,
    error: null,
  })

  const startCountdown = (seconds: number) => {
    if (countdownRef.current) clearInterval(countdownRef.current)
    setState(prev => ({ ...prev, timeRemaining: seconds }))
    countdownRef.current = setInterval(() => {
      setState(prev => {
        if (prev.timeRemaining === null || prev.timeRemaining <= 1) {
          clearInterval(countdownRef.current!)
          countdownRef.current = null
          return { ...prev, timeRemaining: null }
        }
        return { ...prev, timeRemaining: prev.timeRemaining - 1 }
      })
    }, 1000)
  }

  const stopCountdown = () => {
    if (countdownRef.current) {
      clearInterval(countdownRef.current)
      countdownRef.current = null
    }
    setState(prev => ({ ...prev, timeRemaining: null }))
  }

  useEffect(() => {
    if (!userID || !carId) return

    const client = new WsClient()
    clientRef.current = client

    const url = `${WEBSOCKET_URL}${BackendEndpoints.WS_DRIVERS}?userID=${userID}&carID=${carId}&latitude=${latitude}&longitude=${longitude}`

    const unsubs = [
      client.on(WsEventType.DriverRegister, (driver) => {
        setState(prev => ({ ...prev, driver }))
      }),
      client.on(WsEventType.DriverTripRequest, (trip) => {
        setState(prev => ({ ...prev, requestedTrip: trip, tripStatus: WsEventType.DriverTripRequest }))
        startCountdown(15)
      }),
      client.on(WsEventType.DriverTripRequestExpired, () => {
        stopCountdown()
        setState(prev => ({ ...prev, tripStatus: null, requestedTrip: null }))
      }),
      client.on(WsEventType.DriverAssigned, (trip) => {
        setState(prev => ({ ...prev, assignedTrip: trip }))
      }),
    ]

    client.connect(url)

    return () => {
      unsubs.forEach(unsub => unsub())
      stopCountdown()
      client.disconnect()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [userID, carId])

  const sendRaw = (type: string, data: unknown) => {
    clientRef.current?.sendRaw(type, data)
  }

  const reset = () => {
    stopCountdown()
    setState(prev => ({ ...prev, tripStatus: null, requestedTrip: null, assignedTrip: null }))
  }

  return { ...state, sendRaw, reset }
}
