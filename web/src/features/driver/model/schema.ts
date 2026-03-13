import { z } from 'zod';

import { CarPackageSlug } from '@/features/packages';

export const createProfileSchema = z.object({
  name: z
    .string()
    .min(2, 'Name must be at least 2 characters')
    .max(60, 'Name is too long'),
  profile_picture: z
    .string()
    .url('Must be a valid URL (https://...)')
    .optional()
    .or(z.literal('')),
});

export type CreateProfileFormData = z.infer<typeof createProfileSchema>;

export const addCarSchema = z.object({
  car_plate: z
    .string()
    .min(2, 'Car plate is required')
    .max(20, 'Car plate is too long')
    .transform((v) => v.trim().toUpperCase()),
  package_slug: z.enum(
    [CarPackageSlug.SEDAN, CarPackageSlug.SUV, CarPackageSlug.VAN, CarPackageSlug.LUXURY],
    { error: 'Select a car type' },
  ),
});

export type AddCarFormData = z.infer<typeof addCarSchema>;
