package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/iamtbay/tyr-fintech/internal/dto"
	"github.com/iamtbay/tyr-fintech/internal/models"
	"github.com/iamtbay/tyr-fintech/internal/services"
	"github.com/iamtbay/tyr-fintech/internal/worker"
)

type mockTransactionRepository struct {
	transferFunc                  func(ctx context.Context, tx *dto.TransferRequest) error
	getUserFunc                   func(ctx context.Context, userID string) ([]models.Transaction, error)
	getTransactionsByWalletIDFunc func(ctx context.Context, walletID string) ([]*models.Transaction, error)
}

func (m *mockTransactionRepository) Transfer(ctx context.Context, tx *dto.TransferRequest) error {
	return m.transferFunc(ctx, tx)
}
func (m *mockTransactionRepository) GetTransactionsByWalletID(ctx context.Context, walletID string) ([]*models.Transaction, error) {
	return m.getTransactionsByWalletIDFunc(ctx, walletID)
}

func TestTransactionService_Transfer(t *testing.T) {
	tests := []struct {
		name               string
		inputFromWalletNum int64
		inputToWalletNum   int64
		inputAmount        int64
		inputTransactionID string
		mockCreateErr      error
		wantErr            bool
	}{
		{
			name:               "success",
			inputFromWalletNum: 1000000001,
			inputToWalletNum:   1000000002,
			inputTransactionID: "tx-1",
			inputAmount:        100,
			mockCreateErr:      nil,
			wantErr:            false,
		},
		{
			name:               "empty idempotency key",
			inputFromWalletNum: 1000000001,
			inputToWalletNum:   1000000002,
			inputTransactionID: "",
			inputAmount:        100,
			mockCreateErr:      errors.New("idempotency key is required"),
			wantErr:            true,
		},
		{
			name:               "database error",
			inputFromWalletNum: 1000000001,
			inputToWalletNum:   1000000002,
			inputTransactionID: "tx-1",
			inputAmount:        100,
			mockCreateErr:      errors.New("insufficent balance"),
			wantErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockTransactionRepository{
				transferFunc: func(ctx context.Context, tx *dto.TransferRequest) error {
					return tt.mockCreateErr
				},
			}
			service := services.NewTransactionService(mockRepo)
			err := service.Transfer(context.Background(), &dto.TransferRequest{
				TransactionID:    tt.inputTransactionID,
				FromWalletNumber: tt.inputFromWalletNum,
				ToWalletNumber:   tt.inputToWalletNum,
				Amount:           tt.inputAmount,
			})

			if tt.name == "success" {
				select {
				case event := <-worker.WebHookQueue:
					if event.TransactionID != tt.inputTransactionID {
						t.Errorf("Expected transaction id %s, but got %s", tt.inputTransactionID, event.TransactionID)
					}
					if event.Amount != tt.inputAmount {
						t.Errorf("Expected amount %d, but got %d", tt.inputAmount, event.Amount)
					}
				default:
					t.Errorf("After succesfully process, no webhook event sent.")
				}
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Transfer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransactionService_GetHistory(t *testing.T) {
}
