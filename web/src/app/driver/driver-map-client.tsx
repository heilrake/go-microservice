'use client'

import dynamic from "next/dynamic"

const DriverMap = dynamic(() => import("@/features/driver/ui/driver-map").then(mod => mod.DriverMap), { ssr: false })

export function DriverMapClient({ carId, userID }: { carId: string; userID: string }) {
  return <DriverMap carId={carId} userID={userID} />
}
