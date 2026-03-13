import { NextResponse } from 'next/server';

import { gatewayFetch } from '@/lib/api-gateway';
import { setAuthCookie } from '@/lib/cookie';
import type { OAuthProviderType } from '@/shared/libs/contracts';

type OAuthRequestBody = {
  code: string;
  provider: OAuthProviderType;
  role: 'driver' | 'rider';
  redirect_uri?: string;
};

export async function POST(request: Request) {
  let body: OAuthRequestBody;
  try {
    body = await request.json() as OAuthRequestBody;
  } catch {
    return NextResponse.json({ message: 'Invalid JSON' }, { status: 400 });
  }

  const { code, provider, role, redirect_uri } = body;

  if (!code || !provider || !role) {
    return NextResponse.json({ message: 'Missing code, provider or role' }, { status: 400 });
  }

  const res = await gatewayFetch('/auth/oauth', {
    method: 'POST',
    body: JSON.stringify({ code, provider, role, redirect_uri }),
  });

  const data = await res.json().catch(() => ({})) as Record<string, unknown>;

  if (!res.ok) {
    return NextResponse.json(
      { message: (data as { message?: string }).message ?? 'OAuth login failed' },
      { status: res.status },
    );
  }

  const token = (data?.data as { token?: string })?.token;
  // Return user/driver info without the token — token goes into the cookie
  const response = NextResponse.json({ data: { ...(data?.data as object ?? {}), token: undefined } });

  if (token) setAuthCookie(response, token);
  return response;
}
