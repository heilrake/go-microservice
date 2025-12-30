export type Trip = {
    id: string;
    userID: string;
    status: string;
    selectedFare: RouteFare;
    route: Route;
    driver?: Driver;
    trip: Trip;
}

export type RequestRideProps = {
    pickup: [number, number],
    destination: [number, number],
}

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

export const CarPackageSlug = {
    SEDAN: "sedan",
    SUV: "suv",
    VAN: "van",
    LUXURY: "luxury",
} as const;

export type CarPackageSlugType = typeof CarPackageSlug[keyof typeof CarPackageSlug];

export type RouteFare = {
    id: string,
    packageSlug: CarPackageSlugType,
    basePrice: number,
    totalPriceInCents?: number,
    expiresAt: Date,
    route: Route,
}


export type HTTPTripStartResponse = {
    tripID: string;
}

export type TripPreview = {
    tripID: string,
    route: [number, number][],
    rideFares: RouteFare[],
    duration: number,
    distance: number,
}


export type Driver = {
    id: string;
    location: Coordinate;
    geohash: string;
    name: string;
    profilePicture: string;
    carPlate: string;
}
