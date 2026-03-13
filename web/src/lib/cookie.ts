import type { NextResponse } from 'next/server';

export const AUTH_COOKIE = 'auth_token';
const IS_PROD = process.env.NODE_ENV === 'production';

/** Set the HttpOnly auth cookie on a NextResponse. */
export function setAuthCookie(res: NextResponse, token: string): NextResponse {
  res.cookies.set(AUTH_COOKIE, token, {
    httpOnly: true,
    secure: IS_PROD,
    sameSite: 'lax',
    path: '/',
    maxAge: 60 * 60 * 24, // 24 h
  });
  return res;
}

/** Clear the auth cookie (logout). */
export function clearAuthCookie(res: NextResponse): NextResponse {
  res.cookies.set(AUTH_COOKIE, '', {
    httpOnly: true,
    secure: IS_PROD,
    sameSite: 'lax',
    path: '/',
    maxAge: 0,
  });
  return res;
}

/** Decode JWT payload without verification (safe — used only server-side). */
export function decodeTokenPayload(token: string): { user_id: string; role: string } | null {
  try {
    const payload = token.split('.')[1];
    return JSON.parse(Buffer.from(payload, 'base64url').toString()) as {
      user_id: string;
      role: string;
    };
  } catch {
    return null;
  }
}
