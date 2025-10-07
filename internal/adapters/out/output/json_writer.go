package output

import (
	"context"
	"encoding/json"
	"io"
	"mf-statement/internal/domain"
)

type Writer interface {
	Write(ctx context.Context, s domain.Statement) error
}

type JSONWriter struct{ W io.Writer }

func NewJSON(w io.Writer) *JSONWriter { return &JSONWriter{W: w} }

func (j *JSONWriter) Write(ctx context.Context, s domain.Statement) error {
	enc := json.NewEncoder(j.W)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}
