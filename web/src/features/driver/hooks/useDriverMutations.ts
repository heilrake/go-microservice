'use client';

import { mutate } from 'swr';

import type { Car, HTTPCreateCarRequest, HTTPCreateDriverRequest } from '@/shared/libs/contracts';

import { swrKeys } from '@/lib/swr/swr-keys';

export function useDriverMutations() {
  const createProfile = async (payload: HTTPCreateDriverRequest) => {
    const res = await fetch(swrKeys.driver.root(), {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });
    if (!res.ok) throw new Error('Failed to create profile');
    const newDriver = await res.json();
    mutate(swrKeys.driver.root(), newDriver.data, false);
  };

  const addCar = async (payload: HTTPCreateCarRequest) => {
    const res = await fetch(swrKeys.driver.cars(), {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });
    if (!res.ok) throw new Error('Failed to add car');
    const newCar = await res.json();
    mutate(swrKeys.driver.cars(), (prev: Car[] | undefined) => [...(prev ?? []), newCar.data], false);
  };

  return { createProfile, addCar };
}