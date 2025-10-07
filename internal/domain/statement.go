package domain

import (
	"fmt"
)

type Statement struct {
	Period           string           `json:"period"`
	TotalIncome      int64            `json:"total_income"`
	TotalExpenditure int64            `json:"total_expenditure"`
	Transactions     []TransactionDTO `json:"transactions"`
}

func NewStatement(period string, transactions []Transaction, totalIncome, totalExpenditure int64) Statement {
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
		Transactions:     transactionDTOs,
	}
}
