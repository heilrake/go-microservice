import type { ParamRoute,StaticRoute } from './types';

export const routes = {
  driver: {
    root: (() => '/driver') satisfies StaticRoute<'/driver'>,

    profile: (() => '/driver/profile') satisfies StaticRoute<
      '/driver/profile'
    >,

    car: ((id: string) =>
      `/driver/cars/${id}`) satisfies ParamRoute<
      '/driver/cars',
      string
    >,
  },

  rider: {
    root: (() => '/rider') satisfies StaticRoute<'/rider'>,
  },
};