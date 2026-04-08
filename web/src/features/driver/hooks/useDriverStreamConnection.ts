import { useEffect, useRef, useState } from 'react';

import type { Trip } from "@/features/trip";

import { WEBSOCKET_URL } from "@/shared/libs/constants";
import type { ClientWsMessage } from '@/shared/libs/contracts';
import { BackendEndpoints, TripEvents } from '@/shared/libs/contracts';

import type { Driver } from "../models/types";

type useDriverConnectionProps = {
  location: {
    latitude: number;
    longitude: number;
  };
  geohash: string;
  carId: string;
  userID: string;
}

export const useDriverStreamConnection = ({
  location,
  geohash,
  carId,
  userID,
}: useDriverConnectionProps) => {

  const [requestedTrip, setRequestedTrip] = useState<Trip | null>(null)
  const [tripStatus, setTripStatus] = useState<typeof TripEvents[keyof typeof TripEvents] | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [driver, setDriver] = useState<Driver | null>(null);
  const [timeRemaining, setTimeRemaining] = useState<number | null>(null);
  const countdownRef = useRef<NodeJS.Timeout | null>(null);
  const tripStatusRef = useRef(tripStatus);

  useEffect(() => {

    const websocket = new WebSocket(`${WEBSOCKET_URL}${BackendEndpoints.WS_DRIVERS}?userID=${userID}&carID=${carId}&latitude=${location.latitude}&longitude=${location.longitude}`);
    setWs(websocket);

    websocket.onopen = () => {
      if (location) {
        websocket.send(JSON.stringify({
          type: TripEvents.DriverLocation,
          data: { location, geohash },
        }));
      }
    };

    websocket.onmessage = (event) => {
      const message = JSON.parse(event.data) as { type: string; data: any };

      switch (message.type) {
        case TripEvents.DriverTripRequest: {
          const trip = message.data?.trip ?? message.data;
          setRequestedTrip(trip);
          setTripStatus(TripEvents.DriverTripRequest);
          startCountdown(15);
          break;
        }
        case TripEvents.DriverTripRequestExpired:
          resetTripStatus();
          break;
        case TripEvents.DriverRegister:
          setDriver(message.data);
          break;
        case 'driver.cmd.error':
          setError(`Registration failed: ${message.data as string}`);
          break;
        default:
          console.warn(`Unhandled WS message type: "${message.type}"`);
      }
    };

    websocket.onclose = () => { console.log('WebSocket closed'); };
    websocket.onerror = (event) => {
      setError('WebSocket error occurred');
      console.error('WebSocket error:', event);
    };

    return () => {
      if (websocket.readyState === WebSocket.OPEN) websocket.close();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [userID, carId]);

  const startCountdown = (seconds: number) => {
    if (countdownRef.current) clearInterval(countdownRef.current);
    setTimeRemaining(seconds);
    countdownRef.current = setInterval(() => {
      setTimeRemaining((prev) => {
        if (prev === null || prev <= 1) {
          clearInterval(countdownRef.current!);
          countdownRef.current = null;
          return null;
        }
        return prev - 1;
      });
    }, 1000);
  };

  const stopCountdown = () => {
    if (countdownRef.current) {
      clearInterval(countdownRef.current);
      countdownRef.current = null;
    }
    setTimeRemaining(null);
  };

  const sendMessage = (message: ClientWsMessage) => {
    if (ws?.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(message));
    } else {
      setError('WebSocket is not connected');
    }
  };

  const resetTripStatus = () => {
    stopCountdown();
    setTripStatus(null);
    setRequestedTrip(null);
  };

  return { error, tripStatus, driver, requestedTrip, timeRemaining, resetTripStatus, sendMessage, setTripStatus };
};
