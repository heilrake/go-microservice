'use client';

import Link from 'next/link';

import type { DriverProfile } from '@/shared/libs/contracts';
import { Button } from '@/shared/ui/button';

import { useDriverData } from '../hooks/useDriverData';
import { useDriverMutations } from '../hooks/useDriverMutations';
import { CarsList } from './cars-list';
import { DriverProfileCard } from './driver-profile-card';
import { routes } from '@/lib/routes/routes';

export function DriverProfilePage() {
  const { driver, cars, isLoading } = useDriverData();
  console.log("cars", cars
    
  )
  const { createProfile, addCar } = useDriverMutations();
  
  const profileExists = driver !== null && driver !== undefined;

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-gray-500 text-sm animate-pulse">Loading profile...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white border-b sticky top-0 z-10">
        <div className="max-w-lg mx-auto px-4 py-4 flex items-center justify-between">
          <Link
            href={routes.driver.root()}
            className="text-sm text-gray-500 hover:text-gray-800 flex items-center gap-1"
          >
            ← Back to Drive
          </Link>
          <h1 className="text-base font-semibold">Driver Profile</h1>
          <div className="w-20" />
        </div>
      </header>

      <div className="max-w-lg mx-auto px-4 py-6 space-y-6">
        <DriverProfileCard
          profile={driver as DriverProfile | null}
          profileExists={profileExists}
          onCreateProfile={createProfile}
        />

        {driver && <CarsList cars={cars} onAddCar={addCar} />}

        {profileExists && cars && cars.length > 0 && (
          <Link href={routes.driver.root()}>
            <Button className="w-full py-6 text-base">Go Online →</Button>
          </Link>
        )}
      </div>
    </div>
  );
}
