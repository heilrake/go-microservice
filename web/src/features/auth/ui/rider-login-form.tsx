"use client"

import { useForm } from "react-hook-form"
import { useRouter } from "next/navigation"
import { zodResolver } from "@hookform/resolvers/zod"

import { API_URL } from "@/shared/libs/constants"
import type { HTTPUserLoginRequestPayload, HTTPUserLoginResponse } from "@/shared/libs/contracts"
import { BackendEndpoints } from "@/shared/libs/contracts"
import { Button } from "@/shared/ui/button"
import { Input } from "@/shared/ui/input"

import { type LoginFormData, loginSchema } from "../model/schema"

export function RiderLoginForm() {
  const router = useRouter()

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    setError,
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  })
  
  const onSubmit = async (data: LoginFormData) => {
    try {
      const payload: HTTPUserLoginRequestPayload = data

      const response = await fetch(`${API_URL}${BackendEndpoints.RIDER_LOGIN}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      })

      if (!response.ok) {
        let message = `Помилка: ${response.status}`
        try {
          const errorData = await response.json()
          message = errorData.message || message
        } catch {
          // не вдалося розпарсити → залишиться статус
        }
        throw new Error(message)
      }

      const json = (await response.json()) as { data: HTTPUserLoginResponse }
      const { user, token } = json.data

      localStorage.setItem("userID", user.id)
      if (token) {
        localStorage.setItem("authToken", token)
      }

      router.push("/ride")
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Не вдалося увійти"

      setError("root", { message })
    }
  }
  return (
    <div className="min-h-screen bg-gradient-to-b from-white to-gray-50 flex items-center justify-center px-4">
      <div className="bg-white p-8 rounded-2xl shadow-lg max-w-md w-full">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">Rider Login</h1>
          <p className="text-gray-600">Sign in to start riding</p>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          <div>
            <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
              Email
            </label>
            <Input
              id="email"
              type="email"
              autoComplete="email"
              placeholder="your@email.com"
              {...register("email")}
              // className={errors.email ? "border-red-500" : ""}
            />
            {errors.email && (
              <p className="mt-1 text-sm text-red-600">{errors.email.message}</p>
            )}
          </div>

          <div>
            <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-2">
              Password
            </label>
            <Input
              id="password"
              type="password"
              autoComplete="current-password"
              placeholder="••••••••"
              {...register("password")}
              // className={errors.password ? "border-red-500" : ""}
            />
            {errors.password && (
              <p className="mt-1 text-sm text-red-600">{errors.password.message}</p>
            )}
          </div>

          {errors.root && (
            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg text-sm">
              {errors.root.message}
            </div>
          )}

          <Button
            type="submit"
            className="w-full text-lg py-6"
            disabled={isSubmitting}
          >
            {isSubmitting ? "Signing in..." : "Sign In"}
          </Button>
        </form>

        <div className="mt-6 text-center">
          <p className="text-sm text-gray-600">
            Don&apos;t have an account?{" "}
            <button
              type="button"
              onClick={() => router.push("/auth/user/register")}
              className="text-primary hover:underline font-medium"
            >
              Sign up
            </button>
          </p>
        </div>

        <div className="mt-4 text-center">
          <button
            type="button"
            onClick={() => router.push("/")}
            className="text-sm text-gray-500 hover:text-gray-700"
          >
            ← Back to home
          </button>
        </div>
      </div>
    </div>
  )
}
