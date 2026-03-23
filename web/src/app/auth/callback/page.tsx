'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';

import { routes } from '@/lib/routes/routes';

export default function OAuthCallbackPage() {
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const code = params.get('code');
    const provider = params.get('provider');
    const role = params.get('role') as 'driver' | 'rider' | null;

    if (!code || !provider || !role) {
      setError('Missing OAuth parameters');
      return;
    }

    const redirectUri = `${window.location.origin}/auth/callback?provider=${provider}&role=${role}`;

    fetch('/api/auth/oauth', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ code, provider, role, redirect_uri: redirectUri }),
    })
      .then(async (res) => {
        if (!res.ok) {
          const data = await res.json().catch(() => ({})) as { message?: string };
          throw new Error(data.message ?? `OAuth login failed (${res.status})`);
        }
        const data = await res.json() as { data?: { user?: { username?: string } } };
        const name = data?.data?.user?.username ?? '';

        if (role === 'rider') {
          router.replace(routes.rider.root());
          return;
        }

        // For driver: check if driver profile already exists
        const driverRes = await fetch('/api/driver');
        if (driverRes.ok) {
          router.replace(routes.driver.profile());
        } else {
          // New driver — go to onboarding
          const query = name ? `?name=${encodeURIComponent(name)}` : '';
          router.replace(`${routes.driver.onboarding()}${query}`);
        }
      })
      .catch((err: unknown) => {
        setError(err instanceof Error ? err.message : 'Unknown error');
      });
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800 flex items-center justify-center px-4">
      {!error && (
        <div className="text-center">
          <div className="w-10 h-10 border-2 border-white/30 border-t-white rounded-full animate-spin mx-auto mb-4" />
          <p className="text-white text-lg">Signing you in...</p>
        </div>
      )}
      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-6 py-4 rounded-xl text-sm max-w-sm text-center">
          <p className="font-medium mb-1">Login failed</p>
          <p>{error}</p>
          <a href="/" className="mt-3 inline-block text-xs underline text-red-500">← Back to home</a>
        </div>
      )}
    </div>
  );
}
