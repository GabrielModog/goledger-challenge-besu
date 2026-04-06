package handler

import (
	"context"
	"math/big"
	"time"

	"github.com/gofiber/fiber/v2"
	"goledger-challenge-besu/internal/models"
	"goledger-challenge-besu/internal/repository"
	"goledger-challenge-besu/pkg/blockchain"
)

type StorageHandler struct {
	blockchain *blockchain.Client
	repo       repository.StorageRepository
}

func NewStorageHandler(bc *blockchain.Client, repo repository.StorageRepository) *StorageHandler {
	return &StorageHandler{
		blockchain: bc,
		repo:       repo,
	}
}

func (h *StorageHandler) SetValue(c *fiber.Ctx) error {
	var req models.SetValueRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	if req.Value == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_value",
			Message: "Value is required",
		})
	}

	value, ok := new(big.Int).SetString(req.Value, 10)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "invalid_value",
			Message: "Value must be a valid integer",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	txHash, err := h.blockchain.SetValue(ctx, value)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "blockchain_error",
			Message: err.Error(),
		})
	}

	return c.JSON(models.SetValueResponse{
		TxHash: txHash,
		Value:  req.Value,
	})
}

func (h *StorageHandler) GetValue(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	value, err := h.blockchain.GetValue(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "blockchain_error",
			Message: err.Error(),
		})
	}

	return c.JSON(models.GetValueResponse{
		Value: value.String(),
	})
}

func (h *StorageHandler) SyncValue(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	value, err := h.blockchain.GetValue(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "blockchain_error",
			Message: err.Error(),
		})
	}

	if err := h.repo.SetValue(ctx, value.String()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
	}

	return c.JSON(models.SyncResponse{
		SyncedValue: value.String(),
	})
}

func (h *StorageHandler) CheckValue(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	blockchainValue, err := h.blockchain.GetValue(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "blockchain_error",
			Message: err.Error(),
		})
	}

	dbValue, err := h.repo.GetValue(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
	}

	return c.JSON(models.CheckResponse{
		Equal:      blockchainValue.String() == dbValue,
		Blockchain: blockchainValue.String(),
		Database:   dbValue,
	})
}

func (h *StorageHandler) Health(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	dbStatus := "healthy"
	if _, err := h.repo.GetValue(ctx); err != nil {
		dbStatus = "unhealthy"
	}

	bcStatus := "healthy"
	if err := h.blockchain.Ping(ctx); err != nil {
		bcStatus = "unhealthy"
	}

	status := "healthy"
	statusCode := fiber.StatusOK
	if dbStatus != "healthy" || bcStatus != "healthy" {
		status = "degraded"
		statusCode = fiber.StatusServiceUnavailable
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"status":     status,
		"database":   dbStatus,
		"blockchain": bcStatus,
	})
}
