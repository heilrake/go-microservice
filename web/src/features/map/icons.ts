import L from 'leaflet';

const icons = {
  car: new L.Icon({
    iconUrl: '/icons/car.svg',
    iconSize: [36, 36],
    iconAnchor: [18, 18],
  }),
  pin: new L.Icon({
    iconUrl: '/icons/pin.svg',
    iconSize: [24, 36],
    iconAnchor: [12, 36],
  }),
  user: new L.Icon({
    iconUrl: '/icons/user.svg',
    iconSize: [30, 30],
    iconAnchor: [15, 30],
  }),
} as const;

export type MapIconName = keyof typeof icons;

export function getMapIcon(name: MapIconName): L.Icon {
  return icons[name];
}
