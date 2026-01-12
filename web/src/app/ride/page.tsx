"use client"
import dynamic from "next/dynamic"


const RiderMap = dynamic(() => import("@/features/rider/ui/rider-map"), { ssr: false })

export default function RidePage() {
  return <RiderMap />
}