package repos

import (
	"context"
	"errors"

	"github.com/iamtbay/tyr-fintech/internal/dto"
	"github.com/iamtbay/tyr-fintech/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	pool *pgxpool.Pool
}

func NewTransactionRepository(pool *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{pool: pool}
}

// Transfer
func (r *TransactionRepository) Transfer(ctx context.Context, req *dto.TransferRequest) error {
	if req.FromWalletNumber == req.ToWalletNumber {
		return errors.New("cannot transfer to the same wallet")
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	resIdemp, err := tx.Exec(ctx, `INSERT INTO idempotency_keys(key) VALUES($1) ON CONFLICT (key) DO NOTHING`, req.TransactionID)
	if err != nil {
		return err
	}
	if resIdemp.RowsAffected() == 0 {
		return errors.New("transaction already processed")
	}

	firstNumber := req.FromWalletNumber
	secondNumber := req.ToWalletNumber
	if req.FromWalletNumber > req.ToWalletNumber {
		firstNumber = req.ToWalletNumber
		secondNumber = req.FromWalletNumber
	}

	var firstID, firstCurrency string
	err = tx.QueryRow(ctx, `SELECT id, currency FROM wallets WHERE wallet_number=$1 AND deleted_at IS NULL FOR UPDATE`, firstNumber).Scan(&firstID, &firstCurrency)
	if err != nil {
		return errors.New("wallet not found")
	}

	var secondID, secondCurrency string
	err = tx.QueryRow(ctx, `SELECT id, currency FROM wallets WHERE wallet_number=$1 AND deleted_at IS NULL FOR UPDATE`, secondNumber).Scan(&secondID, &secondCurrency)
	if err != nil {
		return errors.New("wallet not found")
	}

	var fromID, toID, fromCurrency, toCurrency string
	if firstNumber == req.FromWalletNumber {
		fromID = firstID
		fromCurrency = firstCurrency
		toID = secondID
		toCurrency = secondCurrency
	} else {
		fromID = secondID
		fromCurrency = secondCurrency
		toID = firstID
		toCurrency = firstCurrency
	}

	if fromCurrency != toCurrency {
		return errors.New("currency mismatch: cannot transfer between different currencies")
	}

	res, err := tx.Exec(ctx, `UPDATE wallets SET balance=balance-$1 WHERE id=$2 AND balance>=$1`, req.Amount, fromID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("insufficient balance")
	}

	_, err = tx.Exec(ctx, `UPDATE wallets SET balance=balance+$1 WHERE id=$2`, req.Amount, toID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `INSERT INTO transactions(id,from_wallet_id,to_wallet_id,amount,status) VALUES($1,$2,$3,$4,$5)`, req.TransactionID, fromID, toID, req.Amount, models.StatusCompleted)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// GetTransactionsByWalletID
func (r *TransactionRepository) GetTransactionsByWalletID(ctx context.Context, walletID string) ([]*models.Transaction, error) {
	query := `
		SELECT 
		t.id,
		t.from_wallet_id,
		t.to_wallet_id,
		w_from.wallet_number as from_wallet_number,
		w_to.wallet_number as to_wallet_number,
		t.amount,
		t.status,
		t.created_at
		FROM transactions t
		LEFT JOIN wallets w_from ON t.from_wallet_id = w_from.id
		LEFT JOIN wallets w_to ON t.to_wallet_id=w_to.id
		WHERE (t.from_wallet_id = $1 OR t.to_wallet_id=$1) AND t.status='COMPLETED'
		ORDER BY t.created_at DESC;
	`
	rows, err := r.pool.Query(ctx, query, walletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(&t.ID, &t.FromWalletID, &t.ToWalletID, &t.FromWalletNumber, &t.ToWalletNumber, &t.Amount, &t.Status, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}

	return transactions, nil
}
