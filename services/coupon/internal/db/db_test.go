package db

import (
	"arvan-challenge/services/coupon/pkg/env"
	"context"
	"fmt"
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
	dbHandler_test DBHandler
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

	dbHandler_test = NewDBHandler(c, &l)

	fmt.Println("sss")

	code := m.Run()
	os.Exit(code)
}

func TestCreateCoupon(t *testing.T) {
	// Initialize mock data
	id := uuid.New()
	coupon := &Coupon{
		Id:           id,
		Code:         "TESTCODE",
		ExpireDate:   time.Now().AddDate(0, 0, 7),
		ChargeAmount: 50,
		AllowedTimes: 5,
	}

	// Initialize a mock DBHandler instance

	// Test CreateCoupon method
	createdCoupon, err := dbHandler_test.CreateCoupon(coupon, context.Background())

	// Check for errors
	if err != nil {
		t.Errorf("CreateCoupon returned an error: %v", err)
	}

	// Check if coupon is created correctly
	if createdCoupon == nil {
		t.Errorf("CreateCoupon did not return the created coupon")
	}
}

func TestGetCoupon(t *testing.T) {
	// Initialize a mock coupon code
	code := "TESTCODE"

	// Initialize a mock DBHandler instance

	// Test GetCoupon method
	retrievedCoupon, err := dbHandler_test.GetCoupon(code, context.Background())

	// Check for errors
	if err != nil {
		t.Errorf("GetCoupon returned an error: %v", err)
	}

	// Check if coupon is retrieved correctly
	if retrievedCoupon == nil {
		t.Errorf("GetCoupon did not return the retrieved coupon")
	}
}

func TestGetUsersByCoupon(t *testing.T) {
	// Initialize a mock coupon code
	code := "TESTCODE"

	// Initialize a mock DBHandler instance

	// Test GetUsersByCoupon method
	users, err := dbHandler_test.GetUsersByCoupon(code, context.Background())

	// Check for errors
	if err != nil {
		t.Errorf("GetUsersByCoupon returned an error: %v", err)
	}

	// Check if users are retrieved correctly
	if users == nil {
		t.Errorf("GetUsersByCoupon did not return the retrieved users")
	}
}

func TestUseCoupon(t *testing.T) {
	// Initialize mock coupon code and phone number
	code := "TESTCODE"
	phoneNumber := "1234567890"

	// Initialize a mock DBHandler instance

	// Test UseCoupon method
	_, err := dbHandler_test.UseCoupon(code, phoneNumber, context.Background())

	// Check for errors
	if err != nil {
		t.Errorf("UseCoupon returned an error: %v", err)
	}
}

func TestGetOutboxes(t *testing.T) {
	// Initialize a mock DBHandler instance

	// Test GetOutboxes method
	outboxes, err := dbHandler_test.GetOutboxes(context.Background())

	// Check for errors
	if err != nil {
		t.Errorf("GetOutboxes returned an error: %v", err)
	}

	// Check if outboxes are retrieved correctly
	if outboxes == nil {
		t.Errorf("GetOutboxes did not return the retrieved outboxes")
	}
}

func TestDeleteOutbox(t *testing.T) {
	// Initialize mock outbox IDs
	outboxIDs := []int{1, 2, 3}

	// Initialize a mock DBHandler instance

	// Test DeleteOutbox method
	err := dbHandler_test.DeleteOutbox(&outboxIDs, context.Background())

	// Check for errors
	if err != nil {
		t.Errorf("DeleteOutbox returned an error: %v", err)
	}
}
