import type { ParamRoute,StaticRoute } from './types';

export const routes = {
  auth: {
    root: ((role: 'driver' | 'rider') => `/auth?role=${role}`) satisfies ParamRoute<'/auth', 'driver' | 'rider'>,
  },

  driver: {
    root: (() => '/driver') satisfies StaticRoute<'/driver'>,
    profile: (() => '/driver/profile') satisfies StaticRoute<'/driver/profile'>,
    onboarding: (() => '/driver/onboarding') satisfies StaticRoute<'/driver/onboarding'>,
    car: ((id: string) => `/driver/cars/${id}`) satisfies ParamRoute<'/driver/cars', string>,
  },

  rider: {
    root: (() => '/rider') satisfies StaticRoute<'/rider'>,
  },
};