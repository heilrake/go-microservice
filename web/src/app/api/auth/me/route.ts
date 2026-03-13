import { cookies } from 'next/headers';
import { NextResponse } from 'next/server';

import { AUTH_COOKIE, decodeTokenPayload } from '@/lib/cookie';

export async function GET() {
  const store = await cookies();
  const token = store.get(AUTH_COOKIE)?.value;

  if (!token) {
    return NextResponse.json({ message: 'Not authenticated' }, { status: 401 });
  }

  const payload = decodeTokenPayload(token);
  if (!payload?.user_id) {
    return NextResponse.json({ message: 'Invalid token' }, { status: 401 });
  }

  return NextResponse.json({ userID: payload.user_id, role: payload.role });
}
