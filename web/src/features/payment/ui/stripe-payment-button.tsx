"use client"

import { loadStripe } from "@stripe/stripe-js"

import type { PaymentEventSessionCreatedData } from "@/features/payment/models/types"

import { clientEnv } from '@/shared/configs/env';
import { Button } from "@/shared/ui/button"

// Initialize Stripe
const stripePromise = loadStripe(clientEnv.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY!)

type StripePaymentButtonProps = {
  paymentSession: PaymentEventSessionCreatedData
  isLoading?: boolean
}

export const StripePaymentButton = ({
  paymentSession,
  isLoading = false,
}: StripePaymentButtonProps) => {
  const handlePayment = async () => {
    const stripe = await stripePromise

    if (!stripe) {
      console.error("Stripe failed to load")
      return
    }

    const { error } = await stripe.redirectToCheckout({ sessionId: paymentSession.sessionID })

    if (error) {
      console.error("Payment error:", error)
    } 
  }

  if (!clientEnv.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY) {
    return (
      <Button
        disabled
        className="w-full bg-red-500 text-white"
      >
        Stripe API KEY is not set on the NEXTJS app
      </Button>
    )
  }

  return (
    <Button
      onClick={handlePayment}
      disabled={isLoading}
      className="w-full"
    >
      {isLoading ? "Loading..." : `Pay ${paymentSession.amount} ${paymentSession.currency}`}
    </Button>
  )
}

