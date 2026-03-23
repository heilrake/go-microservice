'use client';

import { Suspense } from 'react';
import { useSearchParams } from 'next/navigation';

import { GoogleAuthButton } from '@/features/auth/ui/google-auth-button';

function AuthContent() {
  const searchParams = useSearchParams();
  const role = searchParams.get('role') as 'driver' | 'rider' | null;

  if (role !== 'driver' && role !== 'rider') {
    return (
      <div className="min-h-screen flex items-center justify-center px-4">
        <p className="text-gray-500">Invalid role. Go back and try again.</p>
      </div>
    );
  }

  const isDriver = role === 'driver';

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800 flex items-center justify-center px-4">
      <div className="bg-white rounded-2xl shadow-xl p-8 w-full max-w-sm">
        <div className="text-center mb-8">
          <div className="text-4xl mb-3">{isDriver ? '🚗' : '🗺️'}</div>
          <h1 className="text-2xl font-bold text-gray-900">
            {isDriver ? 'Continue as Driver' : 'Continue as Rider'}
          </h1>
          <p className="text-gray-500 text-sm mt-2">
            {isDriver
              ? 'Sign in to start earning'
              : 'Sign in to book a ride'}
          </p>
        </div>

        <GoogleAuthButton role={role} />

        <p className="text-center text-xs text-gray-400 mt-6">
          New? We'll create your account automatically.
        </p>

        <div className="mt-4 text-center">
          <a href="/" className="text-xs text-gray-400 hover:text-gray-600 underline">
            ← Back to home
          </a>
        </div>
      </div>
    </div>
  );
}

export default function AuthPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800 flex items-center justify-center">
        <div className="animate-pulse text-white">Loading...</div>
      </div>
    }>
      <AuthContent />
    </Suspense>
  );
}
