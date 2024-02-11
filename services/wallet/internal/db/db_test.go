package db

import (
	"arvan-challenge/services/wallet/pkg/env"
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

	dbName := "wallet"
	dbUser := "postgres"
	dbPassword := "postgres"

	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:latest"),
		postgres.WithInitScripts(filepath.Join("../../../../", "wallet_init.sql")),
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

func TestAddTransaction(t *testing.T) {

	tid, _ := uuid.NewV7()
	wid, _ := uuid.NewV7()

	// Initialize a mock wallet and transaction
	wallet := &Wallet{
		Id:           tid,
		PhoneNumber:  "+989109810624",
		LastModefied: time.Now(),
		Inventory:    50,
	}

	transaction := &Transaction{
		Id:          wid,
		Description: "Test transaction",
		Date:        time.Now(),
		Amount:      50,
	}

	// Initialize a new DBHandler instance
	dbHandler := NewDBHandler(cfg, logger)

	// Test AddTransaction method
	updatedWallet, err := dbHandler.AddTransaction(transaction, wallet, context.Background())

	// Check for errors
	if err != nil {
		t.Errorf("AddTransaction returned an error: %v", err)
	}

	// Check if wallet inventory has been updated correctly
	if updatedWallet.Inventory != 50 {
		t.Errorf("Wallet inventory is incorrect, got: %d, want: %d", updatedWallet.Inventory, 50)
	}

	// Close the connection
	dbHandler.CloseConnection()
}

func TestGetUserTransactions(t *testing.T) {

	// Initialize a mock phone number
	phoneNumber := "+989109810624"

	// Initialize a new DBHandler instance
	dbHandler := NewDBHandler(cfg, logger)

	// Test GetUserTransactions method
	transactions, err := dbHandler.GetUserTransactions(phoneNumber, context.Background())

	// Check for errors
	if err != nil {
		t.Errorf("GetUserTransactions returned an error: %v", err)
	}

	// Check if transactions are retrieved correctly
	if len(*transactions) != 1 {
		t.Errorf("Number of transactions retrieved is incorrect, got: %d, want: %d", len(*transactions), 1)
	}

	// Close the connection
	dbHandler.CloseConnection()
}
