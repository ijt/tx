package engine

import (
	"encoding/csv"
	"io"
)

type E struct {
}

func New() *E {
	return &E{}
}

func (e *E) Run(r io.Reader, w io.Writer) error {
	cw := csv.NewWriter(w)
	cw.Write([]string{"client", "available", "held", "total", "locked"})
	return nil
}
