import { cookies } from 'next/headers';

import { AUTH_COOKIE } from './cookie';

const GATEWAY_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8081';

/**
 * Server-side proxy to the api-gateway.
 * Reads the HttpOnly auth cookie and adds it as Authorization: Bearer header.
 */
export async function gatewayFetch(path: string, init: RequestInit = {}): Promise<Response> {
  const store = await cookies();
  const token = store.get(AUTH_COOKIE)?.value;

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(init.headers as Record<string, string> ?? {}),
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  };

  return fetch(`${GATEWAY_URL}${path}`, { ...init, headers });
}
