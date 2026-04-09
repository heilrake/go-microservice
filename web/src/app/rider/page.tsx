"use client"
import dynamic from "next/dynamic"

const RiderMapWidget = dynamic(() => import("@/widgets/rider-map-widget").then(mod => mod.RiderMapWidget), { ssr: false })

export default function RidePage() {
  return <RiderMapWidget />
}