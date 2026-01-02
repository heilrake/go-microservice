import type { Driver } from "@/features/driver/models/types";
import type { CarPackageSlugType } from "@/features/packages";

export type Coordinate = {
  latitude: number,
  longitude: number,
}

export type Route = {
  geometry: {
    coordinates: Coordinate[]
  }[],
  duration: number,
  distance: number,
}

export type RouteFare = {
  id: string,
  packageSlug: CarPackageSlugType,
  basePrice: number,
  totalPriceInCents?: number,
  expiresAt: Date,
  route: Route,
}

export type Trip = {
  id: string;
  userID: string;
  status: string;
  selectedFare: RouteFare;
  route: Route;
  driver?: Driver;
  trip: Trip;
}

export type TripPreview = {
  tripID: string,
  route: [number, number][],
  rideFares: RouteFare[],
  duration: number,
  distance: number,
}

export type RequestRideProps = {
  pickup: [number, number],
  destination: [number, number],
}

export type HTTPTripStartResponse = {
  tripID: string;
}

