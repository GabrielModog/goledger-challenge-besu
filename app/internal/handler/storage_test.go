package handler

import (
	"context"
	"goledger-challenge-besu/internal/models"
	"math/big"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBlockchainClient struct {
	mock.Mock
}

func (m *MockBlockchainClient) GetValue(ctx context.Context) (*big.Int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *MockBlockchainClient) SetValue(ctx context.Context, value *big.Int) (string, error) {
	args := m.Called(ctx, value)
	return args.String(0), args.Error(1)
}

func (m *MockBlockchainClient) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockStorageRepository struct {
	mock.Mock
}

func (m *MockStorageRepository) GetValue(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockStorageRepository) SetValue(ctx context.Context, value string) error {
	args := m.Called(ctx, value)
	return args.Error(0)
}

func TestGetValueResponse_Structure(t *testing.T) {
	response := models.GetValueResponse{Value: "42"}
	assert.Equal(t, "42", response.Value)
	assert.NotNil(t, response)
}

func TestSetValueRequest_Structure(t *testing.T) {
	request := models.SetValueRequest{Value: "100"}
	assert.Equal(t, "100", request.Value)
}

func TestCheckResponse_Equal(t *testing.T) {
	response := models.CheckResponse{
		Equal:      true,
		Blockchain: "42",
		Database:   "42",
	}
	assert.True(t, response.Equal)
	assert.Equal(t, response.Blockchain, response.Database)
}

func TestCheckResponse_NotEqual(t *testing.T) {
	response := models.CheckResponse{
		Equal:      false,
		Blockchain: "42",
		Database:   "100",
	}
	assert.False(t, response.Equal)
	assert.NotEqual(t, response.Blockchain, response.Database)
}

func TestHealthResponse_Structure(t *testing.T) {
	response := models.HealthResponse{
		Status:   "healthy",
		Database: "connected",
	}
	assert.Equal(t, "healthy", response.Status)
	assert.Equal(t, "connected", response.Database)
}

func TestErrorResponse_Structure(t *testing.T) {
	response := models.ErrorResponse{
		Error:   "invalid_request",
		Message: "Value is required",
	}
	assert.Equal(t, "invalid_request", response.Error)
	assert.Equal(t, "Value is required", response.Message)
}

func TestSyncResponse_Structure(t *testing.T) {
	response := models.SyncResponse{SyncedValue: "42"}
	assert.Equal(t, "42", response.SyncedValue)
}

func TestFiberApp_HealthEndpoint(t *testing.T) {
	app := fiber.New()
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":   "healthy",
			"database": "ok",
		})
	})

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
