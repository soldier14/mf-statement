package domain

import (
	"fmt"
	"time"
)

type Statement struct {
	Period           string           `json:"period"`
	TotalIncome      int64            `json:"total_income"`
	TotalExpenditure int64            `json:"total_expenditure"`
	NetAmount        int64            `json:"net_amount"`
	TransactionCount int              `json:"transaction_count"`
	Transactions     []TransactionDTO `json:"transactions"`
	GeneratedAt      time.Time        `json:"generated_at"`
}

func NewStatement(period string, transactions []Transaction, totalIncome, totalExpenditure int64) Statement {
	netAmount := totalIncome + totalExpenditure
	transactionDTOs := make([]TransactionDTO, len(transactions))
	for i, tx := range transactions {
		transactionDTOs[i] = TransactionDTO{
			Date:    tx.Date.Format(CSVDateLayout),
			Amount:  fmt.Sprintf("%d", tx.Amount),
			Content: tx.Content,
		}
	}

	return Statement{
		Period:           period,
		TotalIncome:      totalIncome,
		TotalExpenditure: totalExpenditure,
		NetAmount:        netAmount,
		TransactionCount: len(transactions),
		Transactions:     transactionDTOs,
		GeneratedAt:      time.Now(),
	}
}
