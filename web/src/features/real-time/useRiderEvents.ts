'use client'

import { useEffect, useRef, useState } from 'react'

import type { Driver } from '@/features/driver/models/types'
import type { PaymentEventSessionCreatedData } from '@/features/payment'
import type { Trip } from '@/features/trip'

import { WEBSOCKET_URL } from '@/shared/libs/constants'
import { BackendEndpoints } from '@/shared/libs/contracts'
import { WsClient, WsEventType } from '@/shared/ws'

type RiderEventsState = {
  drivers: Driver[]
  tripStatus: WsEventType | null
  assignedDriver: Trip['driver'] | null
  paymentSession: PaymentEventSessionCreatedData | null
}

export function useRiderEvents(userID: string) {
  const clientRef = useRef<WsClient | null>(null)
  const [state, setState] = useState<RiderEventsState>({
    drivers: [],
    tripStatus: null,
    assignedDriver: null,
    paymentSession: null,
  })

  useEffect(() => {
    if (!userID) return

    const client = new WsClient()
    clientRef.current = client

    const unsubs = [
      client.on(WsEventType.DriverLocation, (drivers) => {
        setState(prev => ({ ...prev, drivers }))
      }),
      client.on(WsEventType.TripCreated, () => {
        setState(prev => ({ ...prev, tripStatus: WsEventType.TripCreated }))
      }),
      client.on(WsEventType.DriverAssigned, (trip) => {
        setState(prev => ({
          ...prev,
          assignedDriver: trip.driver,
          tripStatus: WsEventType.DriverAssigned,
        }))
      }),
      client.on(WsEventType.NoDriversFound, () => {
        setState(prev => ({ ...prev, tripStatus: WsEventType.NoDriversFound }))
      }),
      client.on(WsEventType.TripCompleted, () => {
        setState(prev => ({ ...prev, tripStatus: WsEventType.TripCompleted }))
      }),
      client.on(WsEventType.PaymentSessionCreated, (session) => {
        setState(prev => ({
          ...prev,
          paymentSession: session,
          tripStatus: WsEventType.PaymentSessionCreated,
        }))
      }),
    ]

    client.connect(`${WEBSOCKET_URL}${BackendEndpoints.WS_RIDERS}?userID=${userID}`)

    return () => {
      unsubs.forEach(unsub => unsub())
      client.disconnect()
    }
  }, [userID])

  const reset = () => {
    setState(prev => ({ ...prev, tripStatus: null, paymentSession: null }))
  }

  const sendPaymentConfirmed = (tripID: string) => {
    clientRef.current?.sendRaw('rider.cmd.payment_confirmed', { tripID })
  }

  return { ...state, reset, sendPaymentConfirmed }
}
