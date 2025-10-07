package in

import (
	"context"
	"io"
	"mf-statement/internal/domain"
	"mf-statement/internal/usecase"
	"net/url"
	"os"
)

type CSVFileSource struct{}

func NewCSVFileSource() *CSVFileSource {
	return &CSVFileSource{}
}

func (s *CSVFileSource) Open(ctx context.Context, uri string) (io.ReadCloser, error) {
	if u, err := url.Parse(uri); err == nil && u.Scheme == "file" {
		return os.Open(u.Path)
	}
	return os.Open(uri)
}

type CSVReaderService struct {
	Source usecase.Source
	Parser usecase.Parser
}

func NewCSVReaderService(source usecase.Source, parser usecase.Parser) *CSVReaderService {
	return &CSVReaderService{
		Source: source,
		Parser: parser,
	}
}

func (s *CSVReaderService) ReadTransactions(ctx context.Context, csvFileURI string) ([]domain.Transaction, error) {
	csvReader, err := s.Source.Open(ctx, csvFileURI)
	if err != nil {
		return nil, domain.NewIOError("failed to open CSV source", err)
	}
	defer csvReader.Close()

	transactions, err := s.Parser.Parse(ctx, csvReader)
	if err != nil {
		return nil, domain.NewParseError("failed to parse CSV", err)
	}

	return transactions, nil
}
