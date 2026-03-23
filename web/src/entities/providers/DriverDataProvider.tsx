'use client';

import { SWRConfig } from 'swr';

import type { Driver } from '@/features/driver/models/types';

import type { Car } from '@/shared/libs/contracts';

import { swrKeys } from '@/lib/swr/swr-keys';

export function DriverDataProvider({
  initialDriver,
  initialCars,
  children,
}: {
  initialDriver: Driver | null;
  initialCars: Car[];
  children: React.ReactNode;
}) {
  return (
    <SWRConfig
      value={{
        fallback: {
          [swrKeys.driver.root()]: initialDriver,
          [swrKeys.driver.cars()]: initialCars,
        },
        fetcher: (resource, init) =>
          fetch(resource, init).then(res => res.json()).then((r: { data?: unknown }) => r?.data ?? r),
        revalidateOnMount: false,
        revalidateOnFocus: false,
        revalidateIfStale: false,
      }}>
      {children}
    </SWRConfig>
  );
}
