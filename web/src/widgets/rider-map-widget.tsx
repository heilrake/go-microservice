'use client'

import { useRef } from 'react'
import { MapContainer, Marker, Popup, Rectangle, TileLayer } from 'react-leaflet'
import L from 'leaflet'
import useSWR from 'swr'

import { getGeohashBounds, getMapIcon, MapClickHandler, RoutingControl } from '@/features/map'
import { useRiderEvents } from '@/features/real-time'
import { RiderTripOverview } from '@/features/rider/ui/rider-trip-overview'
import { useTripBooking } from '@/features/trip-booking'

import { authApi } from '@/shared/api'
import { LogoutButton } from '@/shared/ui/logout-button'

const LOCATION = { latitude: 49.438280, longitude: 32.060711 }

const userMarker = getMapIcon('pin')
const driverMarker = getMapIcon('car')

export function RiderMapWidget() {
  const mapRef = useRef<L.Map>(null)

  const { data: me } = useSWR('auth/me', () => authApi.me())
  const userID = me?.userID ?? ''

  const events = useRiderEvents(userID)
  const booking = useTripBooking(userID, LOCATION)

  return (
    <div className="relative flex flex-col md:flex-row h-screen">
      <div className={`${booking.destination ? 'flex-[0.7]' : 'flex-1'}`}>
        <MapContainer
          center={[LOCATION.latitude, LOCATION.longitude]}
          zoom={13}
          style={{ height: '100%', width: '100%' }}
          ref={mapRef}
        >
          <TileLayer
            url="https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png"
            attribution="&copy; OpenStreetMap contributors &copy; CARTO"
          />

          <Marker position={[LOCATION.latitude, LOCATION.longitude]} icon={userMarker} />

          {events.drivers.map((driver) => (
            <Rectangle
              key={`grid-${driver.geohash}`}
              bounds={getGeohashBounds(driver.geohash) as L.LatLngBoundsExpression}
              pathOptions={{ color: '#3388ff', weight: 1, fillOpacity: 0.1 }}
            >
              <Popup>Geohash: {driver.geohash}</Popup>
            </Rectangle>
          ))}

          {events.drivers.map((driver) => (
            <Marker
              key={driver.id}
              position={[driver.location.latitude, driver.location.longitude]}
              icon={driverMarker}
            >
              <Popup>Driver: {driver.name}</Popup>
            </Marker>
          ))}

          {booking.destination && (
            <Marker position={booking.destination} icon={userMarker}>
              <Popup>Destination</Popup>
            </Marker>
          )}

          {booking.trip && <RoutingControl route={booking.trip.route} />}

          <MapClickHandler onClick={(e) => booking.handleMapClick(e.latlng.lat, e.latlng.lng)} />
        </MapContainer>
      </div>

      <div className="flex-[0.4] relative">
        <div className="absolute top-3 right-3 z-10">
          <LogoutButton className="text-xs bg-white/90 hover:bg-white border border-gray-200 rounded-lg px-3 py-1.5 text-gray-600 hover:text-gray-900 shadow-sm transition-colors" />
        </div>
        <RiderTripOverview
          trip={booking.trip}
          assignedDriver={events.assignedDriver}
          status={events.tripStatus}
          paymentSession={events.paymentSession}
          selectedCarType={booking.selectedCarType}
          onPackageSelect={booking.startTrip}
          onCancel={events.reset}
        />
      </div>
    </div>
  )
}
