'use client'

import { useState } from 'react'

import { CardElement, Elements, useElements, useStripe } from '@stripe/react-stripe-js'
import { loadStripe } from '@stripe/stripe-js'

import { clientEnv } from '@/shared/configs/env'
import { Button } from '@/shared/ui/button'

const stripePromise = loadStripe(clientEnv.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY!)

type PaymentFormInnerProps = {
  clientSecret: string
  onConfirmed: () => void
}

function PaymentFormInner({ clientSecret, onConfirmed }: PaymentFormInnerProps) {
  const stripe = useStripe()
  const elements = useElements()
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!stripe || !elements) return

    setLoading(true)
    setError(null)

    const card = elements.getElement(CardElement)
    if (!card) {
      setLoading(false)
      return
    }

    const { error: stripeError, paymentIntent } = await stripe.confirmCardPayment(clientSecret, {
      payment_method: { card },
    })

    setLoading(false)

    if (stripeError) {
      setError(stripeError.message ?? 'Payment failed')
      return
    }

    if (paymentIntent?.status === 'requires_capture') {
      onConfirmed()
    }
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4">
      <div className="border rounded-lg p-3 bg-white">
        <CardElement options={{ style: { base: { fontSize: '16px' } } }} />
      </div>

      {error && <p className="text-sm text-red-500">{error}</p>}

      <Button type="submit" disabled={!stripe || loading} className="w-full">
        {loading ? 'Processing...' : 'Authorize payment'}
      </Button>
    </form>
  )
}

type StripePaymentFormProps = {
  clientSecret: string
  onConfirmed: () => void
}

export function StripePaymentForm({ clientSecret, onConfirmed }: StripePaymentFormProps) {
  if (!clientEnv.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY) {
    return (
      <Button disabled className="w-full bg-red-500 text-white">
        Stripe API KEY is not set
      </Button>
    )
  }

  return (
    <Elements stripe={stripePromise}>
      <PaymentFormInner clientSecret={clientSecret} onConfirmed={onConfirmed} />
    </Elements>
  )
}
