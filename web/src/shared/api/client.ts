const GATEWAY_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8081'

async function request<T>(baseUrl: string, path: string, init: RequestInit = {}): Promise<T> {
  const res = await fetch(`${baseUrl}${path}`, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...(init.headers as Record<string, string> ?? {}),
    },
  })

  if (!res.ok) {
    const err = await res.json().catch(() => ({})) as { message?: string }
    throw new Error(err.message ?? `Request failed: ${res.status}`)
  }

  const json = await res.json() as { data?: T } | T
  return ('data' in (json as object) ? (json as { data: T }).data : json) as T
}

// For calls that go through Next.js API routes (/api/*) — handles auth cookie server-side
export const internalClient = {
  get: <T>(path: string) =>
    request<T>('', path),
  post: <T>(path: string, body: unknown) =>
    request<T>('', path, { method: 'POST', body: JSON.stringify(body) }),
  put: <T>(path: string, body: unknown) =>
    request<T>('', path, { method: 'PUT', body: JSON.stringify(body) }),
  delete: <T>(path: string) =>
    request<T>('', path, { method: 'DELETE' }),
}

// For calls that go directly to the api-gateway from the browser
export const gatewayClient = {
  get: <T>(path: string) =>
    request<T>(GATEWAY_URL, path),
  post: <T>(path: string, body: unknown) =>
    request<T>(GATEWAY_URL, path, { method: 'POST', body: JSON.stringify(body) }),
}
