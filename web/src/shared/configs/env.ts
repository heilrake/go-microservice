import { z } from 'zod';

// const serverSchema = z.object({
// });

const clientSchema = z.object({
  NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY: z.string().min(1),
});

// export const serverEnv = serverSchema.parse({
// });

export const clientEnv = clientSchema.parse({
  NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY:
    process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY,
});
