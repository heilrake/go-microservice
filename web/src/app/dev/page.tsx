'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';

import { routes } from '@/lib/routes/routes';

type Persona = {
  label: string;
  role: 'driver' | 'rider';
  seed: number;
  emoji: string;
};

const PERSONAS: Persona[] = [
  { label: 'Driver #1', role: 'driver', seed: 1, emoji: '🚗' },
  { label: 'Driver #2', role: 'driver', seed: 2, emoji: '🚙' },
  { label: 'Driver #3', role: 'driver', seed: 3, emoji: '🚕' },
  { label: 'Rider #1',  role: 'rider',  seed: 1, emoji: '🧍' },
  { label: 'Rider #2',  role: 'rider',  seed: 2, emoji: '🧍' },
  { label: 'Rider #3',  role: 'rider',  seed: 3, emoji: '🧍' },
];

export default function DevLoginPage() {
  const router = useRouter();
  const [loading, setLoading] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  if (process.env.NODE_ENV !== 'development') {
    return <p className="p-8 text-gray-500">Not available in production.</p>;
  }

  const login = async (persona: Persona) => {
    const key = `${persona.role}-${persona.seed}`;
    setLoading(key);
    setError(null);

    try {
      const res = await fetch('/api/dev/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ role: persona.role, seed: persona.seed }),
      });

      if (!res.ok) {
        const d = await res.json().catch(() => ({})) as { message?: string };
        throw new Error(d.message ?? 'Login failed');
      }

      if (persona.role === 'driver') {
        // check if driver profile exists
        const driverRes = await fetch('/api/driver');
        if (driverRes.ok) {
          router.push(routes.driver.profile());
        } else {
          const name = `Dev Driver ${persona.seed}`;
          router.push(`${routes.driver.onboarding()}?name=${encodeURIComponent(name)}`);
        }
      } else {
        router.push(routes.rider.root());
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
      setLoading(null);
    }
  };

  const drivers = PERSONAS.filter((p) => p.role === 'driver');
  const riders = PERSONAS.filter((p) => p.role === 'rider');

  return (
    <div className="min-h-screen bg-gray-950 flex items-center justify-center px-4">
      <div className="w-full max-w-md">
        {/* Header */}
        <div className="text-center mb-8">
          <div className="inline-flex items-center gap-2 bg-yellow-400/10 border border-yellow-400/20 rounded-full px-4 py-1.5 mb-4">
            <span className="text-yellow-400 text-xs font-semibold uppercase tracking-wider">Dev Only</span>
          </div>
          <h1 className="text-2xl font-bold text-white">Quick Login</h1>
          <p className="text-gray-500 text-sm mt-1">Skip Google OAuth — use a test account</p>
        </div>

        {/* How it works */}
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-4 mb-6 text-xs text-gray-400 space-y-1">
          <p>• Кожен акаунт зберігається в БД з email <code className="text-gray-300">dev-driver-1@dev.local</code></p>
          <p>• Перший логін — створює юзера, наступні — знаходить існуючого</p>
          <p>• JWT містить реальний <code className="text-gray-300">user_id</code> та <code className="text-gray-300">role</code></p>
          <p>• В продакшені цей endpoint повертає 403</p>
        </div>

        {/* Drivers */}
        <div className="mb-4">
          <p className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">Drivers</p>
          <div className="grid grid-cols-3 gap-2">
            {drivers.map((p) => {
              const key = `${p.role}-${p.seed}`;
              return (
                <button
                  key={key}
                  onClick={() => login(p)}
                  disabled={loading !== null}
                  className="flex flex-col items-center gap-1.5 bg-gray-900 hover:bg-gray-800 border border-gray-700 hover:border-indigo-500 rounded-xl py-4 px-2 transition-all disabled:opacity-50 cursor-pointer"
                >
                  <span className="text-2xl">{p.emoji}</span>
                  <span className="text-xs font-medium text-gray-300">
                    {loading === key ? '...' : p.label}
                  </span>
                  <span className="text-[10px] text-gray-600">dev-driver-{p.seed}@dev.local</span>
                </button>
              );
            })}
          </div>
        </div>

        {/* Riders */}
        <div className="mb-6">
          <p className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">Riders</p>
          <div className="grid grid-cols-3 gap-2">
            {riders.map((p) => {
              const key = `${p.role}-${p.seed}`;
              return (
                <button
                  key={key}
                  onClick={() => login(p)}
                  disabled={loading !== null}
                  className="flex flex-col items-center gap-1.5 bg-gray-900 hover:bg-gray-800 border border-gray-700 hover:border-green-500 rounded-xl py-4 px-2 transition-all disabled:opacity-50 cursor-pointer"
                >
                  <span className="text-2xl">{p.emoji}</span>
                  <span className="text-xs font-medium text-gray-300">
                    {loading === key ? '...' : p.label}
                  </span>
                  <span className="text-[10px] text-gray-600">dev-rider-{p.seed}@dev.local</span>
                </button>
              );
            })}
          </div>
        </div>

        {error && (
          <div className="bg-red-900/30 border border-red-700/50 text-red-400 rounded-xl px-4 py-3 text-sm text-center mb-4">
            {error}
          </div>
        )}

        <p className="text-center text-xs text-gray-600">
          <Link href="/" className="hover:text-gray-400 underline">← Back to home</Link>
        </p>
      </div>
    </div>
  );
}
