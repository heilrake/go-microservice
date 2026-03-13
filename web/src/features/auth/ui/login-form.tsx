'use client';

import { useForm } from 'react-hook-form';
import { useRouter } from 'next/navigation';
import { zodResolver } from '@hookform/resolvers/zod';

import { API_URL } from '@/shared/libs/constants';
import { BackendEndpoints, OAuthProviders, type OAuthProviderType } from '@/shared/libs/contracts';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';

import { type LoginFormData, loginSchema } from '../model/schema';

type Props = {
  role: 'driver' | 'rider';
};

export function LoginForm({ role }: Props) {
  const router = useRouter();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    setError,
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: { email: '', password: '' },
  });

  const handleOAuthLogin = (provider: OAuthProviderType) => {
    const redirectUri = `${window.location.origin}/auth/callback?provider=${provider}&role=${role}`;
    let oauthUrl = '';

    if (provider === OAuthProviders.GOOGLE) {
      const clientId = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID;
      oauthUrl =
        `https://accounts.google.com/o/oauth2/v2/auth?client_id=${clientId}` +
        `&redirect_uri=${encodeURIComponent(redirectUri)}` +
        `&response_type=code&scope=openid email profile&prompt=consent`;
    }

    if (provider === OAuthProviders.FACEBOOK) {
      const clientId = process.env.NEXT_PUBLIC_FACEBOOK_CLIENT_ID;
      oauthUrl =
        `https://www.facebook.com/v17.0/dialog/oauth?client_id=${clientId}` +
        `&redirect_uri=${encodeURIComponent(redirectUri)}&response_type=code&scope=email,public_profile`;
    }

    window.location.href = oauthUrl;
  };

  const onSubmit = async (data: LoginFormData) => {
    try {
      const response = await fetch(
        `${API_URL}${role === 'driver' ? BackendEndpoints.DRIVER_LOGIN : BackendEndpoints.RIDER_LOGIN}`,
        {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(data),
          credentials: 'include',
        },
      );

      if (!response.ok) {
        let message = `Помилка: ${response.status}`;
        try {
          const errorData = await response.json();
          message = errorData.message || message;
        } catch {}
        throw new Error(message);
      }

      router.push(role === 'driver' ? '/driver' : '/rider');
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Не вдалося увійти';
      setError('root', { message });
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-white to-gray-50 flex items-center justify-center px-4">
      <div className="bg-white p-8 rounded-2xl shadow-lg max-w-md w-full">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            {role === 'driver' ? 'Driver Login' : 'Rider Login'}
          </h1>
          <p className="text-gray-600">Sign in to continue</p>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          <Input label="Email" type="email" {...register('email')} error={errors.email?.message} />
          <Input
            label="Password"
            type="password"
            {...register('password')}
            error={errors.password?.message}
          />

          {errors.root && <p className="text-red-600 text-sm">{errors.root.message}</p>}

          <Button type="submit" className="w-full py-3" disabled={isSubmitting}>
            {isSubmitting ? 'Signing in...' : 'Sign In'}
          </Button>
        </form>

        <div className="mt-6 text-center">
          <p className="text-sm text-gray-500 mb-2">Or sign in with</p>
          <div className="flex justify-center gap-4">
            <Button
              type="button"
              onClick={() => handleOAuthLogin('google')}
              className="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-lg">
              Google
            </Button>
            <Button
              type="button"
              onClick={() => handleOAuthLogin('facebook')}
              className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg">
              Facebook
            </Button>
          </div>
        </div>

        <div className="mt-4 text-center text-sm text-gray-500">
          <button type="button" onClick={() => router.push('/')} className="hover:underline">
            ← Back to home
          </button>
        </div>
      </div>
    </div>
  );
}
