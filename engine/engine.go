package engine

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

type E struct {
	clients map[string]*client
}

type client struct {
	available float64
	held      float64
	locked    bool
}

func New() *E {
	return &E{clients: make(map[string]*client)}
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

	// Process all transactions.
	lineNumber := 1
	for {
		tx, err := cr.Read()
		lineNumber++
		if err == io.EOF {
			break
		}
		typ := tx[0]
		id := tx[1]
		amtStr := strings.TrimSpace(tx[3])
		amt, err := strconv.ParseFloat(amtStr, 64)
		if err != nil {
			return fmt.Errorf("invalid amount %q at line %d", amtStr, lineNumber)
		}
		cli := e.clients[id]
		if cli == nil {
			cli = &client{}
			e.clients[id] = cli
		}
		switch typ {
		case "deposit":
			cli.available += amt
		case "withdrawal":
			cli.available -= amt
		default:
			return fmt.Errorf("unrecognized transaction type %q at line %d", typ, lineNumber)
		}
	}

	fmt.Fprintln(w, "client, available, held, total, locked")
	for id, cli := range e.clients {
		total := cli.available + cli.held
		fmt.Fprintf(w, "%s, %.4f, %.4f, %.4f, %v\n", id, cli.available, cli.held, total, cli.locked)
	}

	return nil
}
