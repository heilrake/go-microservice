import { NextResponse } from 'next/server';

export async function POST() {
  return NextResponse.json({ message: 'Email/password login is no longer supported. Use Google OAuth.' }, { status: 410 });
}
