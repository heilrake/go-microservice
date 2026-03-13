"use client"
import dynamic from "next/dynamic"
import Link from "next/link"

import { routes } from "@/lib/routes/routes"

const DriverMap = dynamic(() => import("@/features/driver/ui/driver-map").then(mod => mod.DriverMap), { ssr: false })

export default function DrivePage() {
  return (
    <div className="relative h-screen">
      {/* Profile shortcut */}
      <div className="absolute top-4 right-4 z-[1000]">
        <Link
          href={routes.driver.profile()}
          className="flex items-center gap-2 bg-white rounded-full px-4 py-2 shadow-md text-sm font-medium hover:bg-gray-50 transition-colors border"
        >
          <span>👤</span>
          Profile
        </Link>
      </div>
      <DriverMap packageSlug="sedan" />
    </div>
  )
}
