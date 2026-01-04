package bootstrap

import (
	"log"
	"strings"

	stripeApi "github.com/stripe/stripe-go/v81"
	"ride-sharing/services/payment-service/pkg/types"
)

func InitStripe(cfg *types.PaymentConfig) {
	if cfg.StripeSecretKey == "" {
		log.Fatal("STRIPE_SECRET_KEY is not set")
	}

	if !strings.HasPrefix(cfg.StripeSecretKey, "sk_") {
		log.Fatal("Invalid Stripe secret key (expected sk_)")
	}

	stripeApi.Key = cfg.StripeSecretKey

	log.Println("Stripe initialized successfully")
}
