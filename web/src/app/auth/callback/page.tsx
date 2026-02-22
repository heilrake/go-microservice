"use client"

import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"

export default function OAuthCallbackPage() {
  const router = useRouter()
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const params = new URLSearchParams(window.location.search)
    const code = params.get("code")
    const provider = params.get("provider")
    const role = params.get("role")

    if (!code || !provider || !role) {
      setError("Missing OAuth parameters")
      setLoading(false)
      return
    }

    const redirectUri = `${window.location.origin}/auth/callback?provider=${provider}&role=${role}`

    fetch(`/api/auth/oauth`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ code, provider, role, redirect_uri: redirectUri }),
      credentials: "include",
    })
      .then(async (res) => {
        if (!res.ok) {
          const data = await res.json().catch(() => ({}))
          throw new Error(data.message || `OAuth login failed (${res.status})`)
        }
        router.push(role === "driver" ? "/drive" : "/ride")
      })
      .catch((err: unknown) => {
        setError(err instanceof Error ? err.message : "Unknown error")
      })
      .finally(() => setLoading(false))
  }, [router])

  return (
    <div className="min-h-screen flex items-center justify-center px-4">
      {loading && <p className="text-gray-700 text-lg">Logging in...</p>}
      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg text-sm">
          {error}
        </div>
      )}
    </div>
  )
}
