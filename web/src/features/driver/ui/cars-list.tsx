'use client';

import { useState } from 'react';

import type { CarPackageSlugType } from '@/features/packages';
import { PackagesMeta } from '@/features/packages/ui/packages-meta';
import type { Car, HTTPCreateCarRequest } from '@/shared/libs/contracts';
import { Button } from '@/shared/ui/button';

import { AddCarForm } from './add-car-form';

type CarsListProps = {
  cars: Car[];
  onAddCar: (data: HTTPCreateCarRequest) => Promise<void>;
};

function CarItem({ car }: { car: Car }) {
  const meta = PackagesMeta[car.package_slug as CarPackageSlugType];
  return (
    <div className="flex items-center gap-3 p-3 rounded-xl border bg-gray-50">
      <div className="p-2 bg-white rounded-lg border text-gray-600 shrink-0">
        {meta?.icon ?? '🚗'}
      </div>
      <div className="flex-1 min-w-0">
        <p className="font-mono font-semibold text-gray-900 tracking-wider">
          {car.car_plate.toUpperCase()}
        </p>
        <p className="text-xs text-gray-500 capitalize">{meta?.name ?? car.package_slug}</p>
      </div>
    </div>
  );
}

export function CarsList({ cars, onAddCar }: CarsListProps) {
  const [showForm, setShowForm] = useState(false);

  return (
    <div className="bg-white rounded-2xl shadow-sm border">
      <div className="px-6 py-4 border-b flex items-center justify-between">
        <h3 className="font-semibold text-gray-800">
          My Cars{' '}
          {cars.length > 0 && (
            <span className="text-gray-400 font-normal text-sm">({cars.length})</span>
          )}
        </h3>
        {!showForm && (
          <button
            type="button"
            onClick={() => setShowForm(true)}
            className="text-sm font-medium text-black hover:underline"
          >
            + Add Car
          </button>
        )}
      </div>

      <div className="p-4 space-y-3">
        {cars.length === 0 && !showForm && (
          <div className="text-center py-6">
            <p className="text-gray-400 text-sm mb-3">No cars added yet</p>
            <Button variant="outline" size="sm" onClick={() => setShowForm(true)}>
              Add your first car
            </Button>
          </div>
        )}

        {cars.map((car) => (
          <CarItem key={car.id} car={car} />
        ))}

        {showForm && (
          <AddCarForm
            onSubmit={async (data) => {
              await onAddCar(data);
              setShowForm(false);
            }}
            onCancel={() => setShowForm(false)}
          />
        )}
      </div>
    </div>
  );
}
