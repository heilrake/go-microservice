import type { ParamKey,StaticKey } from './types';

export const swrKeys = {
  rider: {
    root: (() => '/rider') satisfies StaticKey<'/rider'>,
  },

  driver: {
    root: (() => '/api/driver') satisfies StaticKey<'/api/driver'>,
    cars: (() => '/api/driver/cars') satisfies StaticKey<'/api/driver/cars'>,
  },

  car: {
    byId: ((id: string) =>
      id ? ['/api/cars', id] as const : null) satisfies ParamKey<
      '/api/cars',
      string
    >,
  },
};