'use client';

import { mutate } from 'swr';

import { driverApi } from '@/shared/api';
import type { Car, HTTPCreateCarRequest, HTTPCreateDriverRequest } from '@/shared/libs/contracts';

import { swrKeys } from '@/lib/swr/swr-keys';

export function useDriverMutations() {
  const createProfile = async (payload: HTTPCreateDriverRequest) => {
    const newDriver = await driverApi.createProfile(payload);
    mutate(swrKeys.driver.root(), newDriver, false);
  };

  const addCar = async (payload: HTTPCreateCarRequest) => {
    const newCar = await driverApi.addCar(payload);
    mutate(swrKeys.driver.cars(), (prev: Car[] | undefined) => [...(prev ?? []), newCar], false);
  };

  return { createProfile, addCar };
}