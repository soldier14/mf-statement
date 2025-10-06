package output

import (
	"context"
	"encoding/json"
	"io"
	"mf-statement/internal/model"
)

type Writer interface {
	Write(ctx context.Context, s model.Statement) error
}

type JSONWriter struct{ W io.Writer }

func NewJSON(w io.Writer) *JSONWriter { return &JSONWriter{W: w} }

func (j *JSONWriter) Write(ctx context.Context, s model.Statement) error {
	enc := json.NewEncoder(j.W)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}
