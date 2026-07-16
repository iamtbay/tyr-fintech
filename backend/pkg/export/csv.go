package export

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"github.com/iamtbay/tyr-fintech/internal/models"
)

func TransactionsToCSV(transactions []*models.Transaction) ([]byte, error) {
	//create a buffer to write the CSV data to ram
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	headers := []string{
		"Transaction ID", "Sender", "Receiver", "Amount", "Status", "Date",
	}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	for _, tx := range transactions {
		row := []string{
			tx.ID,
			fmt.Sprintf("%d", tx.FromWalletNumber),
			fmt.Sprintf("%d", tx.ToWalletNumber),
			fmt.Sprintf("%d", tx.Amount),
			string(tx.Status),
			tx.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
