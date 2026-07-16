package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamtbay/tyr-fintech/internal/dto"
	"github.com/iamtbay/tyr-fintech/internal/models"
)

type WalletService interface {
	CreateWallet(ctx context.Context, req *dto.CreateWallet) error
	GetByUserID(ctx context.Context, userID string) ([]*models.Wallet, error)
	DeleteWallet(ctx context.Context, userID, walletID string) error
}

type WalletHandler struct {
	walletService WalletService
}

func NewWalletHandler(walletService WalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

// CREATE WALLET
func (h *WalletHandler) Create(c *gin.Context) {
	var req dto.CreateWallet
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.walletService.CreateWallet(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Wallet created successfully"})
}

// GET WALLETS
func (h *WalletHandler) GetWallets(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	wallets, err := h.walletService.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallets": wallets})
}

// DELETE WALLET
func (h *WalletHandler) DeleteWallet(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	walletId := c.Param("walletID")
	if walletId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}
	if err := h.walletService.DeleteWallet(c.Request.Context(), userID, walletId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet deleted successfully"})
}
