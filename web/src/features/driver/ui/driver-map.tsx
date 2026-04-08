'use client';

import { useMemo, useRef, useState } from 'react';
import { MapContainer, Marker, Popup, TileLayer } from 'react-leaflet';
import Link from 'next/link';
import L from 'leaflet';
import * as Geohash from 'ngeohash';

import { useDriverStreamConnection } from '@/features/driver/hooks/useDriverStreamConnection';
import { DriverCard } from '@/features/driver/ui/driver-card';
import { DriverTripOverview } from '@/features/driver/ui/driver-trip-overview';
import { getMapIcon, MapClickHandler, RoutingControl } from '@/features/map';
import type { CarPackageSlugType } from '@/features/packages';
import type { Coordinate } from '@/features/trip';

import { TripEvents } from '@/shared/libs/contracts';

import { routes } from '@/lib/routes/routes';

const START_LOCATION: Coordinate = {
  latitude: 49.43828,
  longitude: 32.060711,
};

const driverMarker = getMapIcon('car');
const startLocationMarker = getMapIcon('user');
const destinationMarker = getMapIcon('pin');

export const DriverMap = ({ carId, userID }: { carId: string; userID: string }) => {
  const mapRef = useRef<L.Map>(null);
  const [riderLocation, setRiderLocation] = useState<Coordinate>(START_LOCATION);

  const driverGeohash = useMemo(
    () => Geohash.encode(riderLocation?.latitude, riderLocation?.longitude, 7),
    [riderLocation?.latitude, riderLocation?.longitude],
  );

  const { error, driver, tripStatus, requestedTrip, timeRemaining, sendMessage, setTripStatus, resetTripStatus } =
    useDriverStreamConnection({
      location: riderLocation,
      geohash: driverGeohash,
      carId,
      userID,
    });

  console.log('driver', driver);

  const handleMapClick = (e: L.LeafletMouseEvent) => {
    setRiderLocation({
      latitude: e.latlng.lat,
      longitude: e.latlng.lng,
    });
  };

  const handleAcceptTrip = () => {
    if (!requestedTrip || !requestedTrip.id || !driver) {
      alert('No trip ID found or driver is not set');
      return;
    }

    sendMessage({
      type: TripEvents.DriverTripAccept,
      data: {
        tripID: requestedTrip.id,
        riderID: requestedTrip.userID,
        driver: driver,
      },
    });

    resetTripStatus();
    setTripStatus(TripEvents.DriverTripAccept);
  };

  const handleDeclineTrip = () => {
    if (!requestedTrip || !requestedTrip.id || !driver) {
      alert('No trip ID found or driver is not set');
      return;
    }

    sendMessage({
      type: TripEvents.DriverTripDecline,
      data: {
        tripID: requestedTrip.id,
        riderID: requestedTrip.userID,
        driver: driver,
      },
    });

    setTripStatus(TripEvents.DriverTripDecline);
    resetTripStatus();
  };

  const parsedRoute = useMemo(
    () =>
      requestedTrip?.route?.geometry[0]?.coordinates.map(
        (coord) => [coord?.longitude, coord?.latitude] as [number, number],
      ),
    [requestedTrip],
  );

  // destination is the last coordinate in the route
  const destination = useMemo(
    () =>
      requestedTrip?.route?.geometry[0]?.coordinates[
        requestedTrip?.route?.geometry[0]?.coordinates?.length - 1
      ],
    [requestedTrip],
  );
  // start location is the first coordinate in the route
  const startLocation = useMemo(
    () => requestedTrip?.route?.geometry[0]?.coordinates[0],
    [requestedTrip],
  );

  if (error) {
    return <div>Error: {error}</div>;
  }

  console.log("requestedTrip", requestedTrip)
  console.log("tripStatus",tripStatus)

  return (
    <div className="relative flex flex-col md:flex-row h-screen">
      <div className="flex-1">
        <MapContainer
          center={[riderLocation.latitude, riderLocation.longitude]}
          zoom={13}
          style={{ height: '100%', width: '100%' }}
          ref={mapRef}>
          <TileLayer
            url="https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png"
            attribution="&copy; <a href='https://www.openstreetmap.org/copyright'>OpenStreetMap</a> contributors &copy; <a href='https://carto.com/'>CARTO</a>"
          />

          <Marker position={[riderLocation.latitude, riderLocation.longitude]} icon={driverMarker}>
            <Popup>Geohash: {driverGeohash}</Popup>
          </Marker>

          {startLocation && (
            <Marker
              position={[startLocation.longitude, startLocation.latitude]}
              icon={startLocationMarker}>
              <Popup>Start Location</Popup>
            </Marker>
          )}

          {destination && (
            <Marker
              position={[destination.longitude, destination.latitude]}
              icon={destinationMarker}>
              <Popup>Destination</Popup>
            </Marker>
          )}

          {parsedRoute && <RoutingControl route={parsedRoute} />}

          <MapClickHandler onClick={handleMapClick} />
        </MapContainer>
      </div>

      <div className="flex flex-col md:w-[400px] bg-white border-t md:border-t-0 md:border-l">
        <div className="p-4 border-b">
          <DriverCard driver={driver} packageSlug={(driver?.packageSlug ?? 'sedan') as CarPackageSlugType} />
        </div>
        <div className="flex-1 overflow-y-auto">
          <DriverTripOverview
            trip={requestedTrip}
            status={tripStatus}
            timeRemaining={timeRemaining}
            onAcceptTrip={handleAcceptTrip}
            onDeclineTrip={handleDeclineTrip}
          />
        </div>
      </div>
    </div>
  );
};
