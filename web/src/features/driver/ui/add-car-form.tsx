'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';

import { CarPackageSlug, type CarPackageSlugType } from '@/features/packages';
import { PackagesMeta } from '@/features/packages/ui/packages-meta';
import type { HTTPCreateCarRequest } from '@/shared/libs/contracts';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';

import { addCarSchema, type AddCarFormData } from '../model/schema';

type AddCarFormProps = {
  onSubmit: (data: HTTPCreateCarRequest) => Promise<void>;
  onCancel: () => void;
};

export function AddCarForm({ onSubmit, onCancel }: AddCarFormProps) {
  const {
    register,
    handleSubmit,
    watch,
    setValue,
    setError,
    formState: { errors, isSubmitting },
  } = useForm<AddCarFormData>({
    resolver: zodResolver(addCarSchema),
    defaultValues: {
      car_plate: '',
      package_slug: CarPackageSlug.SEDAN,
    },
  });

  const selectedSlug = watch('package_slug');

  const onValid = async (data: AddCarFormData) => {
    try {
      await onSubmit({ car_plate: data.car_plate, package_slug: data.package_slug });
    } catch (err) {
      setError('root', { message: err instanceof Error ? err.message : 'Failed to add car' });
    }
  };

  return (
    <form onSubmit={handleSubmit(onValid)} className="space-y-4 p-4 bg-gray-50 rounded-xl border">
      <h3 className="font-semibold text-gray-800">Add New Car</h3>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Car Plate</label>
        <Input
          placeholder="e.g. AA1234BB"
          className="font-mono tracking-wider uppercase"
          {...register('car_plate', {
            onChange: (e) => {
              e.target.value = e.target.value.toUpperCase();
            },
          })}
        />
        {errors.car_plate && (
          <p className="mt-1 text-xs text-red-600">{errors.car_plate.message}</p>
        )}
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Car Type</label>
        <div className="grid grid-cols-2 gap-2">
          {(Object.keys(PackagesMeta) as CarPackageSlugType[]).map((slug) => {
            const meta = PackagesMeta[slug];
            const isSelected = selectedSlug === slug;
            return (
              <button
                key={slug}
                type="button"
                onClick={() => setValue('package_slug', slug, { shouldValidate: true })}
                className={[
                  'flex items-center gap-2 p-3 rounded-lg border text-left transition-all',
                  isSelected
                    ? 'border-black bg-black text-white'
                    : 'border-gray-200 hover:border-gray-400 bg-white',
                ].join(' ')}
              >
                <span className={isSelected ? 'text-white' : 'text-gray-600'}>{meta.icon}</span>
                <span className="text-sm font-medium">{meta.name}</span>
              </button>
            );
          })}
        </div>
        {errors.package_slug && (
          <p className="mt-1 text-xs text-red-600">{errors.package_slug.message}</p>
        )}
      </div>

      {errors.root && (
        <p className="text-sm text-red-600 bg-red-50 px-3 py-2 rounded-lg">{errors.root.message}</p>
      )}

      <div className="flex gap-2 pt-1">
        <Button type="submit" disabled={isSubmitting} className="flex-1">
          {isSubmitting ? 'Adding...' : 'Add Car'}
        </Button>
        <Button type="button" variant="outline" onClick={onCancel} className="flex-1">
          Cancel
        </Button>
      </div>
    </form>
  );
}
