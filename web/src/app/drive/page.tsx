"use client"
import dynamic from "next/dynamic"

const DriverMap = dynamic(() => import("@/features/driver/ui/driver-map").then(mod => mod.DriverMap), { ssr: false })

export default function DrivePage() {

   // add request to get package slug from driver profile from backend
  return  <DriverMap packageSlug="sedan" />
}