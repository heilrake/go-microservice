import { NextResponse } from 'next/server';

import { gatewayFetch } from '@/lib/api-gateway';
import { setAuthCookie } from '@/lib/cookie';

export async function POST(request: Request) {
  let body: unknown;
  try {
    body = await request.json();
  } catch {
    return NextResponse.json({ message: 'Invalid JSON' }, { status: 400 });
  }

  const res = await gatewayFetch('/rider/login', {
    method: 'POST',
    body: JSON.stringify(body),
  });

  const data = await res.json().catch(() => ({})) as Record<string, unknown>;

  if (!res.ok) {
    return NextResponse.json(
      { message: (data as { message?: string }).message ?? 'Login failed' },
      { status: res.status },
    );
  }

  const token = (data?.data as { token?: string })?.token;
  const response = NextResponse.json({
    data: { user: (data?.data as { user?: unknown })?.user },
  });

  if (token) setAuthCookie(response, token);
  return response;
}
