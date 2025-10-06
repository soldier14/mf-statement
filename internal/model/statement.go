package model

import "time"

const CSVDateLayout = "2006/01/02" // YYYY/MM/DD

type Transaction struct {
	Date    time.Time
	Amount  int64 // in cents
	Content string
}

type TransactionDTO struct {
	Date    string `json:"date"`
	Amount  string `json:"amount"`
	Content string `json:"content"`
}

type Statement struct {
	Period           string           `json:"period"`
	TotalIncome      int64            `json:"total_income"`
	TotalExpenditure int64            `json:"total_expenditure"`
	Transactions     []TransactionDTO `json:"transactions"`
}
