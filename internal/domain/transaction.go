package domain

import (
	"fmt"
	"strings"
	"time"
)

const CSVDateLayout = "2006/01/02"

type Transaction struct {
	Date    time.Time
	Amount  int64
	Content string
}

type TransactionDTO struct {
	Date    string `json:"date"`
	Amount  string `json:"amount"`
	Content string `json:"content"`
}

func NewTransaction(date time.Time, amount int64, content string) (Transaction, error) {
	if content == "" {
		return Transaction{}, fmt.Errorf("content cannot be empty")
	}

	content = strings.TrimSpace(content)
	if content == "" {
		return Transaction{}, fmt.Errorf("content cannot be empty after trimming")
	}

	return Transaction{
		Date:    date,
		Amount:  amount,
		Content: content,
	}, nil
}

func (t Transaction) IsIncome() bool {
	return t.Amount > 0
}

func (t Transaction) IsExpense() bool {
	return t.Amount < 0
}

func (t Transaction) AbsAmount() int64 {
	if t.Amount < 0 {
		return -t.Amount
	}
	return t.Amount
}
