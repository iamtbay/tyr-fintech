package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/iamtbay/tyr-fintech/internal/dto"
	"github.com/iamtbay/tyr-fintech/internal/models"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet *models.Wallet) error
	GetByUserID(ctx context.Context, userID string) ([]*models.Wallet, error)
	GetByID(ctx context.Context, walletID string) (int64, error)
	Delete(ctx context.Context, userID, walletID string) error
}

type WalletService struct {
	walletRepo WalletRepository
}

func NewWalletService(walletRepo WalletRepository) *WalletService {
	return &WalletService{walletRepo: walletRepo}
}

func (s *WalletService) GetByUserID(ctx context.Context, userID string) ([]*models.Wallet, error) {
	return s.walletRepo.GetByUserID(ctx, userID)
}

// CREATE WALLET
func (s *WalletService) CreateWallet(ctx context.Context, req *dto.CreateWallet) error {
	if req.Currency != "TRY" && req.Currency != "USD" && req.Currency != "EUR" {
		return errors.New("Invalid Currency")
	}
	err := s.walletRepo.Create(ctx, &models.Wallet{
		ID:       uuid.New().String(),
		UserID:   req.UserID,
		Currency: req.Currency,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *WalletService) DeleteWallet(ctx context.Context, userID string, walletID string) error {
	balance, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		return err
	}
	if balance > 0 {
		return errors.New("cannot delete with balance greater than 0")
	}
	return s.walletRepo.Delete(ctx, userID, walletID)
}
