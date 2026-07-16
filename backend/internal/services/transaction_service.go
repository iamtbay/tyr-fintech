package services

import (
	"context"
	"errors"

	"github.com/iamtbay/tyr-fintech/internal/dto"
	"github.com/iamtbay/tyr-fintech/internal/models"
	"github.com/iamtbay/tyr-fintech/internal/worker"
)

type TransactionRepository interface {
	Transfer(ctx context.Context, req *dto.TransferRequest) error
	GetTransactionsByWalletID(ctx context.Context, walletID string) ([]*models.Transaction, error)
}

type TransactionService struct {
	repo TransactionRepository
}

func NewTransactionService(repo TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

// Transfer
func (s *TransactionService) Transfer(ctx context.Context, req *dto.TransferRequest) error {
	if req.TransactionID == "" {
		return errors.New("idempotency key is required")
	}
	err := s.repo.Transfer(ctx, req)
	if err != nil {
		return err
	}
	worker.WebHookQueue <- &dto.TransactionWebhookEvent{
		TransactionID:    req.TransactionID,
		FromWalletNumber: req.FromWalletNumber,
		ToWalletNumber:   req.ToWalletNumber,
		Amount:           req.Amount,
		Status:           "COMPLETED",
	}
	return nil
}

// GetHistory
func (s *TransactionService) GetHistory(ctx context.Context, walletID string) ([]*models.Transaction, error) {
	return s.repo.GetTransactionsByWalletID(ctx, walletID)
}
