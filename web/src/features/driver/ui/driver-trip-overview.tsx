import type { Trip } from "@/features/trip";
import { TripOverviewCard } from "@/features/trip";

import { TripEvents, type TripEventType } from "@/shared/libs/contracts"
import { Button } from "@/shared/ui/button"

type DriverTripOverviewProps = {
  trip?: Trip | null,
  status?: TripEventType | null,
  timeRemaining?: number | null,
  onAcceptTrip?: () => void,
  onDeclineTrip?: () => void
}

export const DriverTripOverview = ({ trip, status, timeRemaining, onAcceptTrip, onDeclineTrip }: DriverTripOverviewProps) => {
  if (!trip) {
    return (
      <TripOverviewCard
        title="Waiting for a rider..."
        description="Waiting for a rider to request a trip..."
      />
    )
  }

  if (status === TripEvents.DriverTripRequest) {
    const TOTAL = 15;
    const remaining = timeRemaining ?? 0;
    const progress = (remaining / TOTAL) * 100;
    const canRespond = remaining > 0;

    return (
      <TripOverviewCard
        title="Trip request received!"
        description="A trip has been requested, check the route and accept the trip if you can take it."
      >
        <div className="flex flex-col gap-3">
          <div className="flex items-center justify-between text-sm">
            <span className="text-gray-500">Time to respond</span>
            <span className={`font-semibold tabular-nums ${remaining <= 5 ? 'text-red-500' : 'text-gray-700'}`}>
              {remaining}s
            </span>
          </div>
          <div className="w-full h-2 bg-gray-200 rounded-full overflow-hidden">
            <div
              className={`h-full rounded-full transition-all duration-1000 ease-linear ${remaining <= 5 ? 'bg-red-500' : 'bg-blue-500'}`}
              style={{ width: `${progress}%` }}
            />
          </div>
          {canRespond ? (
            <div className="flex flex-col gap-2">
              <Button onClick={onAcceptTrip}>Accept trip</Button>
              <Button variant="outline" onClick={onDeclineTrip}>Decline trip</Button>
            </div>
          ) : (
            <p className="text-sm text-center text-gray-400">Time expired — waiting for next assignment</p>
          )}
        </div>
      </TripOverviewCard>
    )
  }

  if (status === TripEvents.DriverTripAccept) {
    return (
      <TripOverviewCard
        title="All set!"
        description="You can now start the trip"
      >
        <div className="flex flex-col gap-4">
          <div className="flex flex-col gap-2">
            <h3 className="text-lg font-bold">Trip details</h3>
            <p className="text-sm text-gray-500">
              Trip ID: {trip.id}
              <br />
              Rider ID: {trip.userID}
            </p>
          </div>
        </div>
      </TripOverviewCard>
    )
  }

  return null
}
