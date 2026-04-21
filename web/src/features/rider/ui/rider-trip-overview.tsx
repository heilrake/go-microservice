import type { Driver } from "@/features/driver/models/types";
import { DriverCard } from "@/features/driver/ui/driver-card";
import type { CarPackageSlugType } from "@/features/packages";
import { StripePaymentForm } from "@/features/payment";
import type { RouteFare, TripPreview } from "@/features/trip";
import { TripOverviewCard } from "@/features/trip";

import { Button } from "@/shared/ui/button"
import { Skeleton } from "@/shared/ui/skeleton"
import { WsEventType } from "@/shared/ws"

import { convertMetersToKilometers, convertSecondsToMinutes } from "../lib/math"
import { DriverList } from "./drivers-list"

type TripOverviewProps = {
  trip: TripPreview | null;
  status: WsEventType | null;
  clientSecret?: string | null;
  assignedDriver?: Driver | null;
  selectedCarType?: CarPackageSlugType | null;
  onPackageSelect: (carPackage: RouteFare) => void;
  onPaymentConfirmed: () => void;
  onCancel: () => void;
}

export const RiderTripOverview = ({
  trip,
  status,
  clientSecret,
  assignedDriver,
  selectedCarType,
  onPackageSelect,
  onPaymentConfirmed,
  onCancel,
}: TripOverviewProps) => {

  if (!trip) {
    return (
      <TripOverviewCard
        title="Start a trip"
        description="Click on the map to set a destination"
      />
    )
  }

  // Payment form — shown immediately after trip created, before driver search
  if (clientSecret && !status) {
    return (
      <TripOverviewCard
        title="Authorize payment"
        description="Enter your card details to hold the fare. You'll only be charged after the trip completes."
      >
        <StripePaymentForm
          clientSecret={clientSecret}
          onConfirmed={onPaymentConfirmed}
        />
      </TripOverviewCard>
    )
  }

  if (status === WsEventType.NoDriversFound) {
    return (
      <TripOverviewCard
        title="No drivers found"
        description="No drivers found for your trip, please try again later"
      >
        <Button variant="outline" className="w-full" onClick={onCancel}>
          Go back
        </Button>
      </TripOverviewCard>
    )
  }

  if (status === WsEventType.TripCompleted) {
    return (
      <TripOverviewCard
        title="Trip completed!"
        description="Your trip is completed, thank you for using our service!"
      >
        <Button variant="outline" className="w-full" onClick={onCancel}>
          Go back
        </Button>
      </TripOverviewCard>
    )
  }

  if (status === WsEventType.DriverAssigned) {
    return (
      <TripOverviewCard
        title="Driver is on the way!"
        description="Your driver has accepted the trip."
      >
        <div className="flex flex-col gap-3">
          <DriverCard driver={assignedDriver} />
          <Button variant="destructive" className="w-full" onClick={onCancel}>
            Cancel trip
          </Button>
        </div>
      </TripOverviewCard>
    )
  }

  if (status === WsEventType.TripCreated) {
    return (
      <TripOverviewCard
        title="Looking for a driver"
        description="Your trip is confirmed! We're matching you with a driver."
      >
        <div className="flex flex-col space-y-3 justify-center items-center mb-4">
          <Skeleton className="h-[125px] w-[250px] rounded-xl" />
          <div className="space-y-2">
            <Skeleton className="h-4 w-[250px]" />
            <Skeleton className="h-4 w-[200px]" />
          </div>
        </div>

        <div className="flex flex-col items-center justify-center gap-2">
          {trip?.duration &&
            <h3 className="text-sm font-medium text-gray-700 mb-2">
              Arriving in: {convertSecondsToMinutes(trip.duration)} · {convertMetersToKilometers(trip.distance ?? 0)}
            </h3>
          }
          <Button variant="destructive" className="w-full" onClick={onCancel}>
            Cancel
          </Button>
        </div>
      </TripOverviewCard>
    )
  }

  if (trip.rideFares && trip.rideFares.length >= 0 && !trip.tripID) {
    return (
      <DriverList
        trip={trip}
        selectedCarType={selectedCarType}
        onPackageSelect={onPackageSelect}
        onCancel={onCancel}
      />
    )
  }

  return null
}
