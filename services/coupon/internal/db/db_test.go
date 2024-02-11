package db

import (
	"arvan-challenge/services/coupon/pkg/env"
	"context"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	logger *zerolog.Logger
	cfg    *env.Config
)

func TestMain(m *testing.M) {

	ctx := context.Background()

	dbName := "coupon"
	dbUser := "postgres"
	dbPassword := "postgres"

	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:latest"),
		postgres.WithInitScripts(filepath.Join("../../../../", "coupon_init.sql")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// Clean up the container
	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	connStr, _ := postgresContainer.Ports(ctx)

	port, _ := strconv.Atoi(connStr["5432/tcp"][0].HostPort)

	// Initialize a mock logger
	l := zerolog.New(os.Stderr)

	// Initialize mock configuration
	c := &env.Config{
		DBUser: dbUser,
		DBPass: dbPassword,
		DBHost: "0.0.0.0",
		DBPort: port,
		DBName: dbName,
	}

	cfg = c
	logger = &l

	code := m.Run()
	os.Exit(code)
}

func TestCreateCoupon(t *testing.T) {
	id, _ := uuid.NewV7()

	// Initialize a mock coupon
	coupon := &Coupon{
		Id:           id,
		Code:         "TESTCODE",
		ExpireDate:   time.Now().AddDate(0, 0, 7), // Expires in 7 days
		ChargeAmount: 50,
		AllowedTimes: 5,
	}

	// Initialize a new DBHandler instance
	dbHandler := NewDBHandler(cfg, logger)

	// Test CreateCoupon method
	createdCoupon, err := dbHandler.CreateCoupon(coupon, nil)

	// Check for errors
	if err != nil {
		t.Errorf("CreateCoupon returned an error: %v", err)
	}

	// Check if coupon is created correctly
	if createdCoupon == nil {
		t.Errorf("CreateCoupon did not return the created coupon")
	}

	// Close the connection
	dbHandler.CloseConnection()
}

func TestGetCoupon(t *testing.T) {

	// Initialize a mock coupon code
	code := "TESTCODE"

	// Initialize a new DBHandler instance
	dbHandler := NewDBHandler(cfg, logger)

	// Test GetCoupon method
	retrievedCoupon, err := dbHandler.GetCoupon(code, context.Background())

	// Check for errors
	if err != nil {
		t.Errorf("GetCoupon returned an error: %v", err)
	}

	// Check if coupon is retrieved correctly
	if retrievedCoupon == nil {
		t.Errorf("GetCoupon did not return the retrieved coupon")
	}

	// Close the connection
	dbHandler.CloseConnection()
}

func TestUseCoupon(t *testing.T) {
	// Initialize mock coupon code and phone number
	code := "TESTCODE"
	phoneNumber := "1234567890"

	// Initialize a new DBHandler instance
	dbHandler := NewDBHandler(cfg, logger)

	// Test UseCoupon method
	_, err := dbHandler.UseCoupon(code, phoneNumber, context.Background())

	// Check for errors
	if err != nil {
		t.Errorf("UseCoupon returned an error: %v", err)
	}

	// Close the connection
	dbHandler.CloseConnection()
}

func TestGetOutboxes(t *testing.T) {

	// Initialize a new DBHandler instance
	dbHandler := NewDBHandler(cfg, logger)

	// Test GetOutboxes method
	outboxes, err := dbHandler.GetOutboxes(context.Background())

	// Check for errors
	if err != nil {
		t.Errorf("GetOutboxes returned an error: %v", err)
	}

	// Check if outboxes are retrieved correctly
	if outboxes == nil {
		t.Errorf("GetOutboxes did not return the retrieved outboxes")
	}

	// Close the connection
	dbHandler.CloseConnection()
}

func TestDeleteOutbox(t *testing.T) {
	// Initialize mock outbox IDs
	outboxIDs := []int{1, 2, 3}

	// Initialize a new DBHandler instance
	dbHandler := NewDBHandler(cfg, logger)

	// Test DeleteOutbox method
	err := dbHandler.DeleteOutbox(&outboxIDs, context.Background())

	// Check for errors
	if err != nil {
		t.Errorf("DeleteOutbox returned an error: %v", err)
	}

	// Close the connection
	dbHandler.CloseConnection()
}
