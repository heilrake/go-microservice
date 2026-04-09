'use client';

import { useRouter } from 'next/navigation';

import { authApi } from '@/shared/api';

export function LogoutButton({ className }: { className?: string }) {
  const router = useRouter();

  const handleLogout = async () => {
    await authApi.logout();
    router.push('/');
  };

  return (
    <button
      onClick={handleLogout}
      className={className ?? 'text-sm text-gray-500 hover:text-gray-800 transition-colors'}
    >
      Log out
    </button>
  );
}
