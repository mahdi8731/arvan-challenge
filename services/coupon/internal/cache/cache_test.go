// package cache_test

// import (
// 	"context"
// 	"log"
// 	"testing"

// 	"github.com/redis/go-redis/v9"
// 	"github.com/testcontainers/testcontainers-go"
// 	"github.com/testcontainers/testcontainers-go/wait"
// )

// func TestWithRedis(t *testing.T) {
// 	ctx := context.Background()
// 	req := testcontainers.ContainerRequest{
// 		Image:        "redis:latest",
// 		ExposedPorts: []string{"6379/tcp"},
// 		WaitingFor:   wait.ForLog("Ready to accept connections"),
// 	}
// 	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
// 		ContainerRequest: req,
// 		Started:          true,
// 	})
// 	if err != nil {
// 		log.Fatalf("Could not start redis: %s", err)
// 	}
// 	defer func() {
// 		if err := redisC.Terminate(ctx); err != nil {
// 			log.Fatalf("Could not stop redis: %s", err)
// 		}
// 	}()

// 	endpoint, err := redisC.Endpoint(ctx, "")
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	client := redis.NewClient(&redis.Options{
// 		Addr: endpoint,
// 	})

// 	_ = client
// }

package cache

import (
	"arvan-challenge/services/coupon/internal/db"
	"arvan-challenge/services/coupon/pkg/env"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestSetKeys(t *testing.T) {
	// Create a new mock Redis client
	mockClient, mock := redismock.NewClientMock()

	// Initialize a mock logger
	logger := zerolog.New(os.Stderr)

	// Initialize mock configuration
	cfg := &env.Config{
		RedisUrl: "http://localhost:6379",
	}

	c := NewCache(cfg, &logger)

	ca, _ := c.(*cache)

	ca.redisClient = mockClient

	// Set up expectations
	mock.ExpectHSet("test_key", map[string]interface{}{"test_field": "test_value"}).SetErr(nil)

	// Execute the method under test
	err := c.SetKeys("test_key", map[string]interface{}{"test_field": "test_value"}, &fasthttp.RequestCtx{})

	// Check assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSetKeyForCoupon(t *testing.T) {
	// Create a new mock Redis client
	mockClient, mock := redismock.NewClientMock()

	// Initialize a mock logger
	logger := zerolog.New(os.Stderr)

	// Create a cache instance with the mock Redis client
	c := &cache{
		redisClient: mockClient,
		l:           &logger,
	}

	// Set up expectations
	mock.ExpectHSet("test_coupon_code", ALLOWED_TIME_KEY, 5).SetErr(nil)
	mock.ExpectExpire("test_coupon_code", time.Until(time.Now().Add(time.Hour))).SetErr(nil)

	// Execute the method under test
	coupon := &db.Coupon{
		Code:         "test_coupon_code",
		AllowedTimes: 5,
		ExpireDate:   time.Now().Add(time.Hour),
	}
	err := c.SetKeyForCoupon(coupon, &fasthttp.RequestCtx{})

	// Check assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFieldExists(t *testing.T) {
	// Create a new mock Redis client
	mockClient, mock := redismock.NewClientMock()

	// Create a cache instance with the mock Redis client
	c := &cache{
		redisClient: mockClient,
	}

	// Set up expectations
	mock.ExpectHExists("test_key", "test_field").SetErr(nil)

	// Execute the method under test
	err := c.FieldExists("test_key", "test_field", &fasthttp.RequestCtx{})

	// Check assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckAndUseCoupon(t *testing.T) {
	// Create a new mock Redis client
	mockClient, mock := redismock.NewClientMock()

	// Create a cache instance with the mock Redis client
	c := &cache{
		redisClient: mockClient,
	}

	// Set up expectations
	// mock.ExpectEvalSha().SetVal("1")

	// Execute the method under test
	err := c.CheckAndUseCoupon("test_key", "test_phone_number", &fasthttp.RequestCtx{})

	// Check assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
