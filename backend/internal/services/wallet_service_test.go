package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/iamtbay/tyr-fintech/internal/models"
	"github.com/iamtbay/tyr-fintech/internal/services"
)

type mockWalletRepository struct {
	funcCreate      func(ctx context.Context, wallet *models.Wallet) error
	funcGetByUserID func(ctx context.Context, userID string) ([]*models.Wallet, error)
	funcGetByID     func(ctx context.Context, walletID string) (int64, error)
	funcDelete      func(ctx context.Context, userID, walletID string) error
}

func (m *mockWalletRepository) Create(ctx context.Context, wallet *models.Wallet) error {
	return m.funcCreate(ctx, wallet)
}
func (m *mockWalletRepository) GetByID(ctx context.Context, walletID string) (int64, error) {
	return m.funcGetByID(ctx, walletID)
}
func (m *mockWalletRepository) GetByUserID(ctx context.Context, userID string) ([]*models.Wallet, error) {
	return m.funcGetByUserID(ctx, userID)
}
func (m *mockWalletRepository) Delete(ctx context.Context, userID, walletID string) error {
	return m.funcDelete(ctx, userID, walletID)
}

// TESTS
func TestWalletService_DeleteWallet(t *testing.T) {
	tests := []struct {
		name          string
		inputUserID   string
		inputWalletID string
		mockBalance   int64
		mockCreateErr error
		wantErr       bool
	}{
		{
			name:          "success",
			inputUserID:   "1",
			inputWalletID: "1",
			mockBalance:   0,
			mockCreateErr: nil,
			wantErr:       false,
		},
		{
			name:          "wallet not found",
			inputUserID:   "1",
			inputWalletID: "1",
			mockCreateErr: errors.New("wallet not found"),
			wantErr:       true,
		},
		{
			name:          "cannot delete with balance",
			inputUserID:   "1",
			inputWalletID: "1",
			mockBalance:   10,
			mockCreateErr: errors.New("cannot delete with balance greater than 0"),
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockWalletRepository{
				funcDelete: func(ctx context.Context, userID, walletID string) error {
					return tt.mockCreateErr
				},
				funcGetByID: func(ctx context.Context, walletID string) (int64, error) {
					if tt.name == "wallet not found" {
						return 0, errors.New("wallet not found")
					}
					return tt.mockBalance, nil
				},
			}
			service := services.NewWalletService(mockRepo)
			err := service.DeleteWallet(context.Background(), tt.inputUserID, tt.inputWalletID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteWallet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}
