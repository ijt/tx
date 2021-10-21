package engine

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strings"
)

type E struct {
}

func New() *E {
	return &E{}
}

func (e *E) Run(r io.Reader, w io.Writer) error {
	cr := csv.NewReader(r)
	cr.TrimLeadingSpace = true
	h, err := cr.Read()
	if err == io.EOF {
		return fmt.Errorf("missing header")
	}
	if err != nil {
		return fmt.Errorf("reading header: %v", err)
	}
	wantHeader := strings.Split("type client tx amount", " ")
	if !reflect.DeepEqual(h, wantHeader) {
		return fmt.Errorf("invalid header: got %#v, want %#v", h, wantHeader)
	}

	fmt.Fprintln(w, "client, available, held, total, locked")

	return nil
}
