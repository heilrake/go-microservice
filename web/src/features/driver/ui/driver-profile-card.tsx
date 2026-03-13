'use client';

import { useState } from 'react';

import type { DriverProfile } from '@/shared/libs/contracts';
import type { HTTPCreateDriverRequest } from '@/shared/libs/contracts';
import { Button } from '@/shared/ui/button';

import { CreateProfileForm } from './create-profile-form';

type DriverProfileCardProps = {
  profile: DriverProfile | null;
  profileExists: boolean;
  onCreateProfile: (data: HTTPCreateDriverRequest) => Promise<void>;
};

function Initials({ name }: { name: string }) {
  const letters = name
    .split(' ')
    .map((n) => n[0])
    .join('')
    .toUpperCase()
    .slice(0, 2);
  return (
    <div className="w-16 h-16 rounded-full bg-black text-white flex items-center justify-center text-xl font-bold shrink-0">
      {letters}
    </div>
  );
}

export function DriverProfileCard({ profile, profileExists, onCreateProfile }: DriverProfileCardProps) {
  const [showForm, setShowForm] = useState(false);

  if (profileExists && profile) {
    return (
      <div className="bg-white rounded-2xl shadow-sm border p-6">
        <div className="flex items-center gap-4">
          {profile.profilePicture ? (
            // eslint-disable-next-line @next/next/no-img-element
            <img
              src={profile.profilePicture}
              alt={profile.name}
              className="w-16 h-16 rounded-full object-cover shrink-0"
            />
          ) : (
            <Initials name={profile.name} />
          )}
          <div>
            <h2 className="text-xl font-bold text-gray-900">{profile.name}</h2>
            <p className="text-sm text-gray-500">Driver profile</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-2xl shadow-sm border p-6">
      <div className="text-center space-y-3">
        <div className="w-16 h-16 rounded-full bg-gray-100 flex items-center justify-center text-2xl mx-auto">
          🚗
        </div>
        <h2 className="text-lg font-semibold text-gray-800">No driver profile yet</h2>
        <p className="text-sm text-gray-500">
          Create your driver profile to start accepting rides
        </p>
        {!showForm && (
          <Button onClick={() => setShowForm(true)} className="w-full">
            Create Profile
          </Button>
        )}
      </div>

      {showForm && (
        <CreateProfileForm
          onSubmit={async (data) => {
            await onCreateProfile(data);
            setShowForm(false);
          }}
          onCancel={() => setShowForm(false)}
        />
      )}
    </div>
  );
}
