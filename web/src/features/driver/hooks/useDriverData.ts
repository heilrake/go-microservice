
import useSWR from 'swr';

import { swrKeys } from '@/lib/swr/swr-keys';

export function useDriverData() {
  const { data: driver, isLoading: driverLoading } = useSWR(swrKeys.driver.root());
  const { data: cars, isLoading: carsLoading } = useSWR(driver ? swrKeys.driver.cars() : null);

  return { driver, cars, isLoading: driverLoading || carsLoading };
}