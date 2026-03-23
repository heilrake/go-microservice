'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';

import type { DriverProfile } from '@/shared/libs/contracts';
import { Button } from '@/shared/ui/button';

import { LogoutButton } from '@/shared/ui/logout-button';

import { useDriverData } from '../hooks/useDriverData';
import { useDriverMutations } from '../hooks/useDriverMutations';
import { CarsList } from './cars-list';
import { DriverProfileCard } from './driver-profile-card';
import { routes } from '@/lib/routes/routes';

export function DriverProfilePage() {
  const router = useRouter();

  const { driver, cars, isLoading } = useDriverData();
  const { createProfile, addCar } = useDriverMutations();

  const [selectedCarId, setSelectedCarId] = useState<string | null>(null);

  const profileExists = driver !== null && driver !== undefined;

  const handleSelectedCarId = (carId: string) => {
    setSelectedCarId(carId);
  };

  const handleToStartDriving = () => {
    if (!selectedCarId) return;
    router.push(`${routes.driver.root()}?carId=${selectedCarId}`);
  };

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
            className="text-sm text-gray-500 hover:text-gray-800 flex items-center gap-1">
            ← Back to Drive
          </Link>
          <h1 className="text-base font-semibold">Driver Profile</h1>
          <LogoutButton />
        </div>
      </header>

      <div className="max-w-lg mx-auto px-4 py-6 space-y-6">
        <DriverProfileCard
          profile={driver as DriverProfile | null}
          profileExists={profileExists}
          onCreateProfile={createProfile}
        />

        {driver && (
          <CarsList
            cars={cars}
            selectedCarId={selectedCarId}
            onAddCar={addCar}
            onSelectCarId={handleSelectedCarId}
          />
        )}

        {profileExists && cars && cars.length > 0 && (
          <Button
            onClick={handleToStartDriving}
            disabled={!selectedCarId}
            className="w-full py-6 text-base">
            Go Online →
          </Button>
        )}
      </div>
    </div>
  );
}
