'use client';

import { Suspense, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';

import { routes } from '@/lib/routes/routes';

const PACKAGE_OPTIONS = [
  { value: 'sedan',  label: 'Sedan',  desc: 'Economic and comfortable' },
  { value: 'suv',    label: 'SUV',    desc: 'Spacious ride for groups' },
  { value: 'van',    label: 'Van',    desc: 'Perfect for larger groups' },
  { value: 'luxury', label: 'Luxury', desc: 'Premium experience' },
];

function OnboardingContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const prefillName = searchParams.get('name') ?? '';

  const [step, setStep] = useState<1 | 2>(1);
  const [name, setName] = useState(prefillName);
  const [carPlate, setCarPlate] = useState('');
  const [packageSlug, setPackageSlug] = useState('sedan');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleStep1 = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    setStep(2);
  };

  const handleStep2 = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!carPlate.trim()) return;
    setLoading(true);
    setError(null);

    try {
      // 1. Create driver profile
      const profileRes = await fetch('/api/driver', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: name.trim() }),
      });
      if (!profileRes.ok) {
        const d = await profileRes.json().catch(() => ({})) as { message?: string };
        throw new Error(d.message ?? 'Failed to create driver profile');
      }

      // 2. Add first car
      const carRes = await fetch('/api/driver/cars', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ car_plate: carPlate.trim(), package_slug: packageSlug }),
      });
      if (!carRes.ok) {
        const d = await carRes.json().catch(() => ({})) as { message?: string };
        throw new Error(d.message ?? 'Failed to add car');
      }

      const carData = await carRes.json() as { data?: { id?: string } };
      const carId = carData?.data?.id;

      router.replace(carId ? `${routes.driver.root()}?carId=${carId}` : routes.driver.profile());
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Something went wrong');
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800 flex items-center justify-center px-4">
      <div className="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
        {/* Progress bar */}
        <div className="h-1 bg-gray-100">
          <div
            className="h-1 bg-indigo-500 transition-all duration-300"
            style={{ width: step === 1 ? '50%' : '100%' }}
          />
        </div>

        <div className="p-8">
          <p className="text-xs font-medium text-indigo-500 uppercase tracking-wide mb-1">
            Step {step} of 2
          </p>

          {step === 1 && (
            <>
              <h1 className="text-2xl font-bold text-gray-900 mb-1">Your profile</h1>
              <p className="text-gray-500 text-sm mb-6">Riders will see this name</p>

              <form onSubmit={handleStep1} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Display name
                  </label>
                  <input
                    type="text"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    placeholder="Your name"
                    required
                    className="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-300"
                  />
                </div>

                <button
                  type="submit"
                  className="w-full bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-xl py-3 transition-colors"
                >
                  Continue →
                </button>
              </form>
            </>
          )}

          {step === 2 && (
            <>
              <h1 className="text-2xl font-bold text-gray-900 mb-1">Add your car</h1>
              <p className="text-gray-500 text-sm mb-6">You can add more cars later</p>

              <form onSubmit={handleStep2} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    License plate
                  </label>
                  <input
                    type="text"
                    value={carPlate}
                    onChange={(e) => setCarPlate(e.target.value.toUpperCase())}
                    placeholder="AA 1234 BB"
                    required
                    className="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-indigo-300"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Car type
                  </label>
                  <div className="space-y-2">
                    {PACKAGE_OPTIONS.map((opt) => (
                      <label
                        key={opt.value}
                        className={`flex items-center gap-3 border rounded-xl px-4 py-3 cursor-pointer transition-colors ${
                          packageSlug === opt.value
                            ? 'border-indigo-400 bg-indigo-50'
                            : 'border-gray-200 hover:border-gray-300'
                        }`}
                      >
                        <input
                          type="radio"
                          name="package"
                          value={opt.value}
                          checked={packageSlug === opt.value}
                          onChange={() => setPackageSlug(opt.value)}
                          className="accent-indigo-600"
                        />
                        <div>
                          <p className="text-sm font-medium text-gray-800">{opt.label}</p>
                          <p className="text-xs text-gray-500">{opt.desc}</p>
                        </div>
                      </label>
                    ))}
                  </div>
                </div>

                {error && (
                  <p className="text-red-600 text-sm">{error}</p>
                )}

                <div className="flex gap-3 pt-1">
                  <button
                    type="button"
                    onClick={() => setStep(1)}
                    className="flex-1 border border-gray-200 text-gray-600 font-medium rounded-xl py-3 hover:bg-gray-50 transition-colors"
                  >
                    ← Back
                  </button>
                  <button
                    type="submit"
                    disabled={loading}
                    className="flex-1 bg-indigo-600 hover:bg-indigo-700 disabled:opacity-60 text-white font-medium rounded-xl py-3 transition-colors"
                  >
                    {loading ? 'Setting up...' : 'Start Driving 🚗'}
                  </button>
                </div>
              </form>
            </>
          )}
        </div>
      </div>
    </div>
  );
}

export default function OnboardingPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800 flex items-center justify-center">
        <div className="animate-pulse text-white">Loading...</div>
      </div>
    }>
      <OnboardingContent />
    </Suspense>
  );
}
