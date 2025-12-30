import { useMapEvents } from 'react-leaflet'

type MapClickHandlerProps = {
  onClick: (e: L.LeafletMouseEvent) => void;
}

export function MapClickHandler({ onClick }: MapClickHandlerProps) {
  useMapEvents({
    click: onClick,
  })
  return null
}

