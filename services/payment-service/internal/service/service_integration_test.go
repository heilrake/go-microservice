package service_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	paymentInternal "ride-sharing/services/payment-service/internal"
	"ride-sharing/services/payment-service/internal/domain"
	"ride-sharing/services/payment-service/internal/infrastructure/repository"
	"ride-sharing/services/payment-service/internal/service"
	sharedBootstrap "ride-sharing/shared/bootstrap"
)

// mockStripe реалізує domain.PaymentProcessor без реального Stripe
type mockStripe struct {
	captureErr error
	cancelErr  error
}

func (m *mockStripe) CreatePaymentIntent(_ context.Context, _ int64, _ string, _ map[string]string) (string, string, error) {
	return "pi_test_123", "secret_test_456", nil
}

func (m *mockStripe) CapturePayment(_ context.Context, _ string, _ *int64) error {
	return m.captureErr
}

func (m *mockStripe) CancelPayment(_ context.Context, _ string) error {
	return m.cancelErr
}

var testDB *gorm.DB

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := tcpostgres.Run(ctx,
		"postgres:16",
		tcpostgres.WithDatabase("payment_test"),
		tcpostgres.WithUsername("test_user"),
		tcpostgres.WithPassword("test_password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
	)
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			log.Printf("failed to terminate container: %v", err)
		}
	}()

	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("failed to get container host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("failed to get container port: %v", err)
	}

	dsn := fmt.Sprintf("postgres://test_user:test_password@%s:%s/payment_test?sslmode=disable", host, port.Port())

	if err := sharedBootstrap.RunMigrator(sharedBootstrap.MigratorConfig{
		MigrationsFS:  paymentInternal.Migrations,
		MigrationsDir: "migrations",
		DatabaseURL:   dsn,
		ServiceName:   "payment-service-test",
	}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect to test db: %v", err)
	}

	os.Exit(m.Run())
}

func newTestService(stripe domain.PaymentProcessor) domain.Service {
	repo := repository.NewPostgresRepository(testDB)
	return service.NewPaymentService(repo, stripe)
}

func cleanDB(t *testing.T) {
	t.Helper()
	testDB.Exec("DELETE FROM payment_intents")
}

func createIntent(t *testing.T, svc domain.Service, tripID string) *domain.PaymentIntentModel {
	t.Helper()
	intent, err := svc.CreatePaymentIntent(context.Background(), tripID, "user-1", 1500, "usd")
	if err != nil {
		t.Fatalf("setup: create intent for trip %s: %v", tripID, err)
	}
	return intent
}

// --- базові сценарії ---

func TestCreatePaymentIntent(t *testing.T) {
	cleanDB(t)
	svc := newTestService(&mockStripe{})

	intent, err := svc.CreatePaymentIntent(context.Background(), "trip-1", "user-1", 1500, "usd")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if intent.StripePaymentIntentID != "pi_test_123" {
		t.Errorf("expected stripe ID pi_test_123, got %s", intent.StripePaymentIntentID)
	}
	if intent.Status != domain.StatusAuthorized {
		t.Errorf("expected status %s, got %s", domain.StatusAuthorized, intent.Status)
	}

	repo := repository.NewPostgresRepository(testDB)
	saved, err := repo.GetByTripID(context.Background(), "trip-1")
	if err != nil {
		t.Fatalf("expected record in db, got error: %v", err)
	}
	if saved.Amount != 1500 {
		t.Errorf("expected amount 1500, got %d", saved.Amount)
	}
}

func TestCapturePayment(t *testing.T) {
	cleanDB(t)
	svc := newTestService(&mockStripe{})
	createIntent(t, svc, "trip-2")

	if err := svc.CapturePayment(context.Background(), "trip-2"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	repo := repository.NewPostgresRepository(testDB)
	saved, _ := repo.GetByTripID(context.Background(), "trip-2")
	if saved.Status != domain.StatusCaptured {
		t.Errorf("expected status %s, got %s", domain.StatusCaptured, saved.Status)
	}
}

func TestCancelPayment(t *testing.T) {
	cleanDB(t)
	svc := newTestService(&mockStripe{})
	createIntent(t, svc, "trip-3")

	if err := svc.CancelPayment(context.Background(), "trip-3"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	repo := repository.NewPostgresRepository(testDB)
	saved, _ := repo.GetByTripID(context.Background(), "trip-3")
	if saved.Status != domain.StatusCancelled {
		t.Errorf("expected status %s, got %s", domain.StatusCancelled, saved.Status)
	}
}

// --- edge cases ---

func TestCapturePayment_TripNotFound(t *testing.T) {
	cleanDB(t)
	svc := newTestService(&mockStripe{})

	err := svc.CapturePayment(context.Background(), "non-existent-trip")
	if err == nil {
		t.Fatal("expected error for non-existent trip, got nil")
	}
}

func TestCancelPayment_TripNotFound(t *testing.T) {
	cleanDB(t)
	svc := newTestService(&mockStripe{})

	err := svc.CancelPayment(context.Background(), "non-existent-trip")
	if err == nil {
		t.Fatal("expected error for non-existent trip, got nil")
	}
}

func TestCreatePaymentIntent_DuplicateTripID(t *testing.T) {
	cleanDB(t)
	svc := newTestService(&mockStripe{})
	createIntent(t, svc, "trip-dup")

	// друге створення для того самого tripID повинно повернути помилку (UNIQUE constraint)
	_, err := svc.CreatePaymentIntent(context.Background(), "trip-dup", "user-1", 1500, "usd")
	if err == nil {
		t.Fatal("expected error for duplicate tripID, got nil")
	}
}

func TestCapturePayment_StripeError_StatusUnchanged(t *testing.T) {
	cleanDB(t)
	svc := newTestService(&mockStripe{captureErr: fmt.Errorf("stripe unavailable")})
	createIntent(t, svc, "trip-stripe-err")

	err := svc.CapturePayment(context.Background(), "trip-stripe-err")
	if err == nil {
		t.Fatal("expected error when stripe fails, got nil")
	}

	repo := repository.NewPostgresRepository(testDB)
	saved, _ := repo.GetByTripID(context.Background(), "trip-stripe-err")
	if saved.Status != domain.StatusAuthorized {
		t.Errorf("expected status to remain %s after stripe error, got %s", domain.StatusAuthorized, saved.Status)
	}
}

// --- concurrency ---

func TestCapturePayment_ConcurrentCalls(t *testing.T) {
	cleanDB(t)
	svc := newTestService(&mockStripe{})
	createIntent(t, svc, "trip-concurrent")

	// 5 горутин намагаються захопити одночасно — тільки одна повинна змінити статус
	const workers = 5
	var wg sync.WaitGroup
	errors := make([]error, workers)

	for i := range workers {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			errors[idx] = svc.CapturePayment(context.Background(), "trip-concurrent")
		}(i)
	}
	wg.Wait()

	// перевіряємо що статус "captured" — незалежно від того скільки горутин завершились успішно
	repo := repository.NewPostgresRepository(testDB)
	saved, _ := repo.GetByTripID(context.Background(), "trip-concurrent")
	if saved.Status != domain.StatusCaptured {
		t.Errorf("expected status %s after concurrent capture, got %s", domain.StatusCaptured, saved.Status)
	}
}

func TestCreatePaymentIntent_ConcurrentDuplicates(t *testing.T) {
	cleanDB(t)
	svc := newTestService(&mockStripe{})

	// 5 горутин намагаються створити intent для одного tripID одночасно
	const workers = 5
	var wg sync.WaitGroup
	results := make([]error, workers)

	for i := range workers {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_, results[idx] = svc.CreatePaymentIntent(context.Background(), "trip-conc-dup", "user-1", 1500, "usd")
		}(i)
	}
	wg.Wait()

	// тільки один успіх — решта помилки через UNIQUE constraint
	successCount := 0
	for _, err := range results {
		if err == nil {
			successCount++
		}
	}
	if successCount != 1 {
		t.Errorf("expected exactly 1 successful create, got %d", successCount)
	}
}
