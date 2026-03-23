import { redirect } from "next/navigation"
import { cookies } from "next/headers"
import Link from "next/link"

import { AUTH_COOKIE, decodeTokenPayload } from "@/lib/cookie"
import { routes } from "@/lib/routes/routes"
import { DriverMapClient } from "./driver-map-client"

export default async function DrivePage({ searchParams }: { searchParams: Promise<{ carId?: string }> }) {
  const { carId } = await searchParams

  if (!carId) redirect(routes.driver.profile())

  const store = await cookies()
  const token = store.get(AUTH_COOKIE)?.value
  const payload = token ? decodeTokenPayload(token) : null

  if (!payload?.user_id) redirect(routes.driver.profile())

  return (
    <div className="relative h-screen">
      <div className="absolute top-4 right-4 z-[1000]">
        <Link
          href={routes.driver.profile()}
          className="flex items-center gap-2 bg-white rounded-full px-4 py-2 shadow-md text-sm font-medium hover:bg-gray-50 transition-colors border"
        >
          <span>👤</span>
          Profile
        </Link>
      </div>
      <DriverMapClient carId={carId} userID={payload.user_id} />
    </div>
  )
}
