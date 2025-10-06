package service

import (
	"context"
	"fmt"
	"io"
	"mf-statement/internal/model"
	"mf-statement/internal/output"
	"mf-statement/internal/parser"
	"sort"
	"strconv"
)

type Source interface {
	Open(context context.Context, uri string) (io.ReadCloser, error)
}

type StatementGenerator struct {
	Source Source
	Parse  parser.Parser
	Write  output.Writer
}

func (generator *StatementGenerator) GenerateMonthly(
	context context.Context,
	periodDisplay string,
	year, month int,
	csvFileURI string,
) error {
	csvReader, err := generator.Source.Open(context, csvFileURI)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer csvReader.Close()

	allTransactions, err := generator.Parse.Parse(context, csvReader)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	monthlyTransactions, totalIncome, totalExpenditure := filterForMonth(allTransactions, year, month)
	sort.SliceStable(monthlyTransactions, func(i, j int) bool { return monthlyTransactions[i].Date.After(monthlyTransactions[j].Date) })

	statement := model.Statement{
		Period:           periodDisplay,
		TotalIncome:      totalIncome,
		TotalExpenditure: totalExpenditure,
		Transactions:     toDTOs(monthlyTransactions),
	}
	return generator.Write.Write(context, statement)
}

func filterForMonth(transactions []model.Transaction, year, month int) ([]model.Transaction, int64, int64) {
	var (
		filteredTransactions []model.Transaction
		totalIncome          int64
		totalExpense         int64
	)
	for _, transaction := range transactions {
		if transaction.Date.Year() == year && int(transaction.Date.Month()) == month {
			filteredTransactions = append(filteredTransactions, transaction)
			if transaction.Amount > 0 {
				totalIncome += transaction.Amount
			} else {
				totalExpense += transaction.Amount
			}
		}
	}
	return filteredTransactions, totalIncome, totalExpense
}

func toDTOs(transactions []model.Transaction) []model.TransactionDTO {
	transactionDTOs := make([]model.TransactionDTO, 0, len(transactions))
	for _, transaction := range transactions {
		transactionDTOs = append(transactionDTOs, model.TransactionDTO{
			Date:    transaction.Date.Format(model.CSVDateLayout),
			Amount:  strconv.FormatInt(transaction.Amount, 10),
			Content: transaction.Content,
		})
	}
	return transactionDTOs
}
