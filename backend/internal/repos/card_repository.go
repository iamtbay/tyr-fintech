package repos

import (
	"context"
	"net/http"

	"github.com/iamtbay/tyr-fintech/internal/models"
	"github.com/iamtbay/tyr-fintech/pkg/apperrors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CardRepository struct {
	db *pgxpool.Pool
}

func NewCardRepository(db *pgxpool.Pool) *CardRepository {
	return &CardRepository{db: db}
}

// CREATE
func (r *CardRepository) Create(ctx context.Context, card *models.Card) error {
	query := `INSERT INTO cards (id,user_id,wallet_id,card_number,cvv,expiry_month,expiry_year,limit_amount, spent_amount,status) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.db.Exec(ctx, query,
		card.ID, card.UserID, card.WalletID, card.CardNumber, card.CVV, card.ExpiryMonth, card.ExpiryYear, card.LimitAmount, card.SpentAmount, card.Status)
	if err != nil {
		return err
	}
	return nil
}

// GET BY USER ID
func (r *CardRepository) GetByUserID(ctx context.Context, userID string) ([]models.Card, error) {
	query := `SELECT id,user_id,wallet_id,card_number,cvv,expiry_month,expiry_year,limit_amount, spent_amount,status,created_at,updated_at FROM cards WHERE user_id = $1`
	cards := []models.Card{}
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		card := models.Card{}
		err := rows.Scan(&card.ID, &card.UserID, &card.WalletID, &card.CardNumber, &card.CVV, &card.ExpiryMonth, &card.ExpiryYear, &card.LimitAmount, &card.SpentAmount, &card.Status, &card.CreatedAt, &card.UpdatedAt)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, nil
}

func (r *CardRepository) UpdateStatus(ctx context.Context, cardID, userID string, status models.CardStatus) error {
	query := `UPDATE cards SET status=$1 WHERE id=$2 AND user_id=$3`
	_, err := r.db.Exec(ctx, query, status, cardID, userID)
	return err
}

//PROCESS PAYMENT

func (r *CardRepository) ProcessPayment(ctx context.Context, cardNumber, cvv string, expiryMonth, expiryYear int, amount int64) error {
	var card models.Card
	query := `SELECT id,wallet_id,limit_amount,spent_amount, status,cvv,expiry_month,expiry_year FROM cards WHERE card_number=$1 FOR UPDATE`
	err := r.db.QueryRow(ctx, query, cardNumber).Scan(&card.ID, &card.WalletID, &card.LimitAmount, &card.SpentAmount, &card.Status, &card.CVV, &card.ExpiryMonth, &card.ExpiryYear)
	if err != nil {
		return err
	}
	if card.Status != models.CardStatusActive {
		return apperrors.New(http.StatusBadRequest, "Card is not active")
	}
	if card.CVV != cvv || card.ExpiryMonth != expiryMonth || card.ExpiryYear != expiryYear {
		return apperrors.New(http.StatusBadRequest, "Invalid card details")
	}
	if card.LimitAmount-card.SpentAmount < int(amount) {
		return apperrors.New(http.StatusBadRequest, "Insufficient funds")
	}

}
