'use client';

import { useRef } from 'react';
import { MapContainer, Marker, Popup, TileLayer } from 'react-leaflet';
import L from 'leaflet';

import { DriverCard } from '@/features/driver/ui/driver-card';
import { DriverTripOverview } from '@/features/driver/ui/driver-trip-overview';
import { useDriverSession } from '@/features/driver-session';
import { getMapIcon, MapClickHandler, RoutingControl } from '@/features/map';
import type { CarPackageSlugType } from '@/features/packages';

const driverMarker = getMapIcon('car');
const startLocationMarker = getMapIcon('user');
const destinationMarker = getMapIcon('pin');

const START_LOCATION = { latitude: 49.43828, longitude: 32.060711 };

type Props = { carId: string; userID: string };

export function DriverMapWidget({ carId, userID }: Props) {
  const mapRef = useRef<L.Map>(null);

  const session = useDriverSession(userID, carId, START_LOCATION);

  const parsedRoute = session.requestedTrip?.route?.geometry[0]?.coordinates.map(
    (coord) => [coord.longitude, coord.latitude] as [number, number],
  );

  const routeCoords = session.requestedTrip?.route?.geometry[0]?.coordinates ?? [];
  const destination = routeCoords[routeCoords.length - 1];
  const startLocation = routeCoords[0];

  return (
    <div className="relative flex flex-col md:flex-row h-screen">
      <div className="flex-1">
        <MapContainer
          center={[session.sessionLocation.latitude, session.sessionLocation.longitude]}
          zoom={13}
          style={{ height: '100%', width: '100%' }}
          ref={mapRef}>
          <TileLayer
            url="https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png"
            attribution="&copy; OpenStreetMap contributors &copy; CARTO"
          />

          <Marker
            position={[session.sessionLocation.latitude, session.sessionLocation.longitude]}
            icon={driverMarker}>
            <Popup>Geohash: {session.geohash}</Popup>
          </Marker>

          {startLocation && (
            <Marker
              position={[startLocation.longitude, startLocation.latitude]}
              icon={startLocationMarker}>
              <Popup>Start</Popup>
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

          <MapClickHandler onClick={(e) => session.moveLocation(e.latlng.lat, e.latlng.lng)} />
        </MapContainer>
      </div>

      <div className="flex flex-col md:w-[400px] bg-white border-t md:border-t-0 md:border-l">
        <div className="p-4 border-b">
          <DriverCard
            driver={session.driver}
            packageSlug={(session.driver?.packageSlug ?? 'sedan') as CarPackageSlugType}
          />
        </div>
        <div className="flex-1 overflow-y-auto">
          <DriverTripOverview
            trip={session.requestedTrip}
            acceptedTrip={session.acceptedTrip}
            status={session.tripStatus}
            timeRemaining={session.timeRemaining}
            onAcceptTrip={() =>
              session.requestedTrip &&
              session.driver &&
              session.acceptTrip(session.requestedTrip, session.driver)
            }
            onDeclineTrip={() =>
              session.requestedTrip &&
              session.driver &&
              session.declineTrip(session.requestedTrip, session.driver)
            }
            onCompleteTrip={session.completeTrip}
          />
        </div>
      </div>
    </div>
  );
}
