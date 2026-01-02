export const CarPackageSlug = {
  SEDAN: "sedan",
  SUV: "suv",
  VAN: "van",
  LUXURY: "luxury",
} as const;

export type CarPackageSlugType = typeof CarPackageSlug[keyof typeof CarPackageSlug];

