import { NextResponse } from 'next/server';

import { gatewayFetch } from '@/lib/api-gateway';

export async function GET() {
  const res = await gatewayFetch('/driver/cars');
  const data = await res.json().catch(() => ({})) as Record<string, unknown>;

  if (res.status === 404) return NextResponse.json({ data: [] }, { status: 200 });
  if (!res.ok) {
    return NextResponse.json(
      { message: (data as { message?: string }).message ?? 'Cars fetch failed' },
      { status: res.status },
    );
  }
  return NextResponse.json(data);
}

export async function POST(request: Request) {
  let body: unknown;
  try {
    body = await request.json();
  } catch {
    return NextResponse.json({ message: 'Invalid JSON' }, { status: 400 });
  }

  const res = await gatewayFetch('/driver/cars', {
    method: 'POST',
    body: JSON.stringify(body),
  });

  const data = await res.json().catch(() => ({})) as Record<string, unknown>;

  if (!res.ok) {
    return NextResponse.json(
      { message: (data as { message?: string }).message ?? 'Add car failed' },
      { status: res.status },
    );
  }
  return NextResponse.json(data);
}
