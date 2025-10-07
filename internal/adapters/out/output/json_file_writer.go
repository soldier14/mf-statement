package output

import (
	"context"
	"mf-statement/internal/domain"
	"os"
)

type JSONFileWriter struct {
	FilePath string
}

func NewJSONFile(filePath string) *JSONFileWriter {
	return &JSONFileWriter{FilePath: filePath}
}

func (j *JSONFileWriter) Write(ctx context.Context, s domain.Statement) error {
	file, err := os.Create(j.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := NewJSON(file)
	return writer.Write(ctx, s)
}
