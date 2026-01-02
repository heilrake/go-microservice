import type { Coordinate } from "@/features/trip";

export type Driver = {
  id: string;
  location: Coordinate;
  geohash: string;
  name: string;
  profilePicture: string;
  carPlate: string;
}
