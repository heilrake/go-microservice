'use server';

import { gatewayFetch } from '@/lib/api-gateway';

export async function getDriverServer() {
  const res = await gatewayFetch('/driver');

  if (res.status === 404) return null;
  if (!res.ok) throw new Error('Driver fetch failed');

  const json = await res.json();
  return json.data;
}

export async function getDriverCarsServer() {
  const res = await gatewayFetch('/driver/cars');

  if (res.status === 404) return [];
  if (!res.ok) throw new Error('Cars fetch failed');

  const json = await res.json();
  return json.data;
}