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

var (
	cHnandler Cache
	mock      redismock.ClientMock
)

func TestMain(m *testing.M) {
	// Create a new mock Redis client
	mockClient, Mock := redismock.NewClientMock()

	// Initialize a mock logger
	logger := zerolog.New(os.Stderr)

	// Initialize mock configuration
	cfg := &env.Config{
		RedisUrl: "http://localhost:6379",
	}

	c := NewCache(cfg, &logger)

	ca, _ := c.(*cache)

	ca.redisClient = mockClient

	cHnandler = c
	mock = Mock

}

func TestSetKeys(t *testing.T) {

	// Set up expectations
	mock.ExpectHSet("test_key", map[string]interface{}{"test_field": "test_value"}).SetErr(nil)

	// Execute the method under test
	err := cHnandler.SetKeys("test_key", map[string]interface{}{"test_field": "test_value"}, &fasthttp.RequestCtx{})

	// Check assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSetKeyForCoupon(t *testing.T) {

	mock.ExpectHSet("test_coupon_code", ALLOWED_TIME_KEY, 5).SetErr(nil)
	mock.ExpectExpire("test_coupon_code", time.Until(time.Now().Add(time.Hour))).SetErr(nil)

	// Execute the method under test
	coupon := &db.Coupon{
		Code:         "test_coupon_code",
		AllowedTimes: 5,
		ExpireDate:   time.Now().Add(time.Hour),
	}
	err := cHnandler.SetKeyForCoupon(coupon, &fasthttp.RequestCtx{})

	// Check assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFieldExists(t *testing.T) {

	// Set up expectations
	mock.ExpectHExists("test_key", "test_field").SetErr(nil)

	// Execute the method under test
	err := cHnandler.FieldExists("test_key", "test_field", &fasthttp.RequestCtx{})

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
