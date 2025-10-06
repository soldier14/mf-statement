package output

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"mf-statement-service/internal/model"
)

type JSONFileWriter struct {
	FilePath string
}

func NewJSONFile(filePath string) *JSONFileWriter {
	return &JSONFileWriter{FilePath: filePath}
}

func (j *JSONFileWriter) Write(ctx context.Context, statement model.Statement) error {
	err := os.MkdirAll(filepath.Dir(j.FilePath), 0755)
	if err != nil {
		return fmt.Errorf("create directories: %w", err)
	}

	file, err := os.Create(j.FilePath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(statement); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	fmt.Printf("Successfully exported statement to %s\n", j.FilePath)
	return nil
}
