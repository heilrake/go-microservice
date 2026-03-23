'use client';

import { useRouter, useSearchParams } from 'next/navigation';
import { Suspense } from 'react';

import { Button } from '@/shared/ui/button';
import { routes } from '@/lib/routes/routes';

function PaymentSuccess() {
  const router = useRouter();
  return (
    <div className="min-h-screen flex items-center justify-center px-4">
      <div className="bg-white p-8 rounded-2xl shadow-lg text-center max-w-md w-full">
        <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg className="w-8 h-8 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <h1 className="text-2xl font-bold text-gray-900">Payment Successful!</h1>
        <p className="text-gray-600 mt-2 mb-6">Your ride has been confirmed.</p>
        <Button className="w-full" variant="outline" onClick={() => router.push('/')}>
          Return Home
        </Button>
      </div>
    </div>
  );
}

function HomeContent() {
  const router = useRouter();
  const searchParams = useSearchParams();

  if (searchParams.get('payment') === 'success') {
    return <PaymentSuccess />;
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800 flex flex-col items-center justify-center px-4">
      <div className="text-center mb-12">
        <h1 className="text-5xl font-bold text-white mb-3">RideShare</h1>
        <p className="text-slate-400 text-lg">Your ride, your way</p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-6 w-full max-w-2xl">
        <button
          onClick={() => router.push(routes.auth.root('rider'))}
          className="group bg-white/5 hover:bg-white/10 border border-white/10 hover:border-white/20 rounded-2xl p-8 text-left transition-all duration-200 cursor-pointer"
        >
          <div className="text-4xl mb-4">🗺️</div>
          <h2 className="text-xl font-semibold text-white mb-1">I want to Ride</h2>
          <p className="text-slate-400 text-sm">Book a trip in seconds</p>
          <div className="mt-6 text-sm font-medium text-indigo-400 group-hover:text-indigo-300">
            Get started →
          </div>
        </button>

        <button
          onClick={() => router.push(routes.auth.root('driver'))}
          className="group bg-white/5 hover:bg-white/10 border border-white/10 hover:border-white/20 rounded-2xl p-8 text-left transition-all duration-200 cursor-pointer"
        >
          <div className="text-4xl mb-4">🚗</div>
          <h2 className="text-xl font-semibold text-white mb-1">I want to Drive</h2>
          <p className="text-slate-400 text-sm">Earn on your schedule</p>
          <div className="mt-6 text-sm font-medium text-indigo-400 group-hover:text-indigo-300">
            Get started →
          </div>
        </button>
      </div>
    </div>
  );
}

export default function Home() {
  return (
    <Suspense fallback={
      <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800 flex items-center justify-center">
        <div className="animate-pulse text-white text-lg">Loading...</div>
      </div>
    }>
      <HomeContent />
    </Suspense>
  );
}
