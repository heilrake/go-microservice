import { useEffect, useState } from 'react';

import type { Driver } from "@/features/driver/models/types";
import type { PaymentEventSessionCreatedData } from "@/features/payment";
import type { Coordinate } from "@/features/trip";
import type { Trip } from "@/features/trip";

import { WEBSOCKET_URL } from "@/shared/libs/constants";
import type { ServerWsMessage } from '@/shared/libs/contracts';
import { BackendEndpoints, isValidWsMessage, TripEvents } from '@/shared/libs/contracts';

export function useRiderStreamConnection(location: Coordinate, userID: string) {
  const [drivers, setDrivers] = useState<Driver[]>([]);
  const [tripStatus, setTripStatus] = useState<typeof TripEvents[keyof typeof TripEvents] | null>(null);
  const [paymentSession, setPaymentSession] = useState<PaymentEventSessionCreatedData | null>(null);
  const [assignedDriver, setAssignedDriver] = useState<Trip["driver"] | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!userID) return;

    const ws = new WebSocket(`${WEBSOCKET_URL}${BackendEndpoints.WS_RIDERS}?userID=${userID}`);

    ws.onopen = () => {
      if (location) {
        ws.send(JSON.stringify({
          type: TripEvents.DriverLocation,
          data: {
            location,
          }
        }));
      }
    };

    ws.onmessage = (event) => {
      const message = JSON.parse(event.data) as ServerWsMessage;

      if (!message || !isValidWsMessage(message)) {
        setError(`Unknown message type "${message}", allowed types are: ${Object.values(TripEvents).join(', ')}`);
        return;
      }

      switch (message.type) {
        case TripEvents.DriverLocation:
          setDrivers(message.data);
          break;
        case TripEvents.PaymentSessionCreated:
          setPaymentSession(message.data);
          setTripStatus(message.type);
          break;
        case TripEvents.DriverAssigned:
          setAssignedDriver(message.data.driver);
          setTripStatus(message.type);
          break;
        case TripEvents.Created:
          setTripStatus(message.type);
          break;
        case TripEvents.NoDriversFound:
          setTripStatus(message.type);
          break;
      }
    };

    ws.onclose = () => {
      console.log('WebSocket closed');
    };

    ws.onerror = (event) => {
      setError('WebSocket error occurred');
      console.error('WebSocket error:', event);
    };

    return () => {
      console.log('Closing WebSocket');
      if (ws.readyState === WebSocket.OPEN) {
        ws.close();
      }
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [userID]);

  const resetTripStatus = () => {
    setTripStatus(null);
    setPaymentSession(null);
  }

  return { drivers, assignedDriver, error, tripStatus, paymentSession, resetTripStatus };
}

