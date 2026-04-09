'use client'

import dynamic from "next/dynamic"

const DriverMapWidget = dynamic(() => import("@/widgets/driver-map-widget").then(mod => mod.DriverMapWidget), { ssr: false })

export function DriverMapClient({ carId, userID }: { carId: string; userID: string }) {
  return <DriverMapWidget carId={carId} userID={userID} />
}
