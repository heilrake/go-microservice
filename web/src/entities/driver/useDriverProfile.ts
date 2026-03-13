'use client';

import useSWR from 'swr';

import { swrKeys } from '@/lib/swr/swr-keys';

const fetcher = (url: string) =>
  fetch(url).then(res => res.json()).then(r => r.data);

export function useDriverProfile() {
  const { data: driver } = useSWR(
    swrKeys.driver.root(),
    fetcher
  );

  const { data: cars } = useSWR(
    driver
      ? swrKeys.driver.cars()
      : null,
    fetcher
  );

  return {
    driver,
    cars: cars ?? [],
    profileExists: !!driver,
  };
}