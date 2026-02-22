import { NextResponse } from "next/server"

import type { OAuthProviderType } from "@/shared/libs/contracts"

type OAuthRequestBody = {
  code: string
  provider: OAuthProviderType
  role: "driver" | "rider"
  redirect_uri?: string
}

export async function POST(request: Request) {
  let body: OAuthRequestBody
  try {
    body = await request.json()
  } catch {
    return NextResponse.json({ message: "Invalid JSON" }, { status: 400 })
  }

  const { code, provider, role, redirect_uri } = body

  if (!code || !provider || !role) {
    return NextResponse.json(
      { message: "Missing code, provider or role" },
      { status: 400 }
    )
  }

  const apiGatewayUrl =
    process.env.NEXT_PUBLIC_API_GATEWAY_URL ||
    process.env.NEXT_PUBLIC_API_URL ||
    "http://localhost:8081"

  try {
    const response = await fetch(`${apiGatewayUrl}/auth/oauth`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ code, provider, role, redirect_uri }),
    })

    const data = await response.json().catch(() => ({}))

    if (!response.ok) {
      return NextResponse.json(
        { message: (data as { message?: string }).message || "OAuth login failed" },
        { status: response.status }
      )
    }

    return NextResponse.json(data)
  } catch (err) {
    console.error("OAuth proxy error:", err)
    return NextResponse.json(
      { message: "Internal server error" },
      { status: 500 }
    )
  }
}
