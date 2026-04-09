import type { Car, DriverProfile,HTTPCreateCarRequest, HTTPCreateDriverRequest } from '@/shared/libs/contracts'

import { internalClient } from './client'

export const driverApi = {
  getProfile: () =>
    internalClient.get<DriverProfile | null>('/api/driver'),

  createProfile: (payload: HTTPCreateDriverRequest) =>
    internalClient.post<DriverProfile>('/api/driver', payload),

  getCars: () =>
    internalClient.get<Car[]>('/api/driver/cars'),

  addCar: (payload: HTTPCreateCarRequest) =>
    internalClient.post<Car>('/api/driver/cars', payload),
}
