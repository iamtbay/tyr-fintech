package export

import (
	"bytes"
	"fmt"
	"time"

	"github.com/iamtbay/tyr-fintech/internal/models"
	"github.com/jung-kurt/gofpdf"
)

func TransactionsToPDF(walletID string, transactions []*models.Transaction) ([]byte, error) {

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	//1 App Header
	pdf.SetFont("Arial", "B", 18)
	pdf.SetTextColor(59, 130, 246)
	pdf.CellFormat(190, 10, "Tyr-Fintech", "0", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(100, 116, 139)
	pdf.CellFormat(190, 8, "Account Statement", "0", 1, "C", false, 0, "")
	pdf.Ln(6)

	//2 Details
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(95, 6, fmt.Sprintf("Wallet ID: %s", walletID), "0", 0, "L", false, 0, "")
	pdf.CellFormat(95, 6, fmt.Sprintf("Generated At: %s", time.Now().Format("2006-01-02 15:04:05")), "0", 1, "R", false, 0, "")
	pdf.Ln(6)

	//3 Table (Grid) Headers
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(241, 245, 249)

	pdf.CellFormat(55, 8, "Transaction ID", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 8, "Sender", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 8, "Receiver", "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 8, "Amount", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Date", "1", 1, "C", true, 0, "")

	//4 TABLE DATA
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(30, 41, 59)
	for _, tx := range transactions {
		amount := tx.Amount
		if tx.FromWalletID != walletID {
			amount = tx.ConvertedAmount
		}
		//showing only first 8 digit of wallet number for security
		pdf.CellFormat(55, 8, tx.ID[:8], "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 8, fmt.Sprintf("%d", tx.FromWalletNumber), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 8, fmt.Sprintf("%d", tx.ToWalletNumber), "1", 0, "C", false, 0, "")
		pdf.CellFormat(25, 8, fmt.Sprintf("%d", amount), "1", 0, "R", false, 0, "")
		pdf.CellFormat(30, 8, tx.CreatedAt.Format("2006-01-02 15:04:05"), "1", 1, "C", false, 0, "")
	}

	// TRANSFER TO RAM
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
