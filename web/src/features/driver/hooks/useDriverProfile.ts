'use client';

import { useCallback, useMemo } from 'react';
import useSWR, { mutate } from 'swr';

import type { Car, DriverProfile, HTTPCreateCarRequest, HTTPCreateDriverRequest } from '@/shared/libs/contracts';

type UseDriverProfileReturn = {
  profile: DriverProfile | null;
  cars: Car[];
  isLoading: boolean;
  error: string | null;
  profileExists: boolean;
  createProfile: (data: HTTPCreateDriverRequest) => Promise<void>;
  addCar: (data: HTTPCreateCarRequest) => Promise<void>;
  refetch: () => void;
};

// Cookies are sent automatically for same-origin requests — no Authorization header needed.
async function fetcher<T>(url: string): Promise<T | null> {
  const res = await fetch(url);
  if (res.status === 401) throw new Error('Unauthorized');
  if (res.status === 404) return null;
  if (!res.ok) throw new Error(`Request failed: ${res.status}`);
  const json = await res.json() as { data: T };
  return json.data;
}

const PROFILE_KEY = '/api/driver';
const CARS_KEY = '/api/driver/cars';

export function useDriverProfile(): UseDriverProfileReturn {
  const {
    data: profile,
    error: profileError,
    isLoading: profileLoading,
  } = useSWR<DriverProfile | null>(PROFILE_KEY, fetcher<DriverProfile>, {
    dedupingInterval: 10_000,
  });

  const profileExists = profile !== null && profile !== undefined;

  // Only fetch cars after we know the profile exists
  const carsKey = useMemo(
    () => (!profileLoading && profileExists ? CARS_KEY : null),
    [profileLoading, profileExists],
  );

  const {
    data: carsData,
    error: carsError,
    isLoading: carsLoading,
  } = useSWR<Car[] | null>(carsKey, fetcher<Car[]>, {
    dedupingInterval: 10_000,
  });

  const refetch = useCallback(() => {
    void mutate(PROFILE_KEY);
    void mutate(CARS_KEY);
  }, []);

  const createProfile = useCallback(async (data: HTTPCreateDriverRequest) => {
    const res = await fetch('/api/driver', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!res.ok) {
      const text = await res.text();
      throw new Error(text || `Error ${res.status}`);
    }
    await mutate(PROFILE_KEY);
  }, []);

  const addCar = useCallback(async (data: HTTPCreateCarRequest) => {
    const res = await fetch('/api/driver/cars', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!res.ok) {
      const text = await res.text();
      throw new Error(text || `Error ${res.status}`);
    }
    await mutate(CARS_KEY);
  }, []);

  const error = profileError?.message ?? carsError?.message ?? null;
  const isLoading = profileLoading || (profileExists && carsLoading);

  return {
    profile: profile ?? null,
    cars: carsData ?? [],
    isLoading,
    error,
    profileExists,
    createProfile,
    addCar,
    refetch,
  };
}
