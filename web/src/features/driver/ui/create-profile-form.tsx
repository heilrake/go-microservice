'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';

import type { HTTPCreateDriverRequest } from '@/shared/libs/contracts';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';

import { createProfileSchema, type CreateProfileFormData } from '../model/schema';

type CreateProfileFormProps = {
  onSubmit: (data: HTTPCreateDriverRequest) => Promise<void>;
  onCancel: () => void;
};

export function CreateProfileForm({ onSubmit, onCancel }: CreateProfileFormProps) {
  const {
    register,
    handleSubmit,
    setError,
    formState: { errors, isSubmitting },
  } = useForm<CreateProfileFormData>({
    resolver: zodResolver(createProfileSchema),
    defaultValues: { name: '', profile_picture: '' },
  });

  const onValid = async (data: CreateProfileFormData) => {
    try {
      await onSubmit({
        name: data.name,
        profile_picture: data.profile_picture || undefined,
      });
    } catch (err) {
      setError('root', {
        message: err instanceof Error ? err.message : 'Failed to create profile',
      });
    }
  };

  return (
    <form onSubmit={handleSubmit(onValid)} className="mt-6 space-y-4">
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Your Name</label>
        <Input placeholder="e.g. John Doe" {...register('name')} />
        {errors.name && (
          <p className="mt-1 text-xs text-red-600">{errors.name.message}</p>
        )}
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          Profile Picture URL{' '}
          <span className="text-gray-400 font-normal">(optional)</span>
        </label>
        <Input placeholder="https://..." {...register('profile_picture')} />
        {errors.profile_picture && (
          <p className="mt-1 text-xs text-red-600">{errors.profile_picture.message}</p>
        )}
      </div>

      {errors.root && (
        <p className="text-sm text-red-600 bg-red-50 px-3 py-2 rounded-lg">
          {errors.root.message}
        </p>
      )}

      <div className="flex gap-2">
        <Button type="submit" disabled={isSubmitting} className="flex-1">
          {isSubmitting ? 'Creating...' : 'Create Profile'}
        </Button>
        <Button type="button" variant="outline" onClick={onCancel} className="flex-1">
          Cancel
        </Button>
      </div>
    </form>
  );
}
