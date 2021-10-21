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
	clients      map[string]*client
	transactions map[string]transaction
}

type client struct {
	available float64
	held      float64
	locked    bool
}

type transaction struct {
	typ    string
	amount float64
}

func New() *E {
	return &E{
		clients:      make(map[string]*client),
		transactions: make(map[string]transaction),
	}
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
		row, err := cr.Read()
		lineNumber++
		if err == io.EOF {
			break
		}
		typ := row[0]
		id := row[1]
		cli := e.clients[id]
		if cli == nil {
			cli = &client{}
			e.clients[id] = cli
		}

		if len(row) == 3 {
			switch typ {
			case "dispute":
				txID := strings.TrimSpace(row[2])
				dtx := e.transactions[txID]
				cli.available -= dtx.amount
				cli.held += dtx.amount

			case "resolve":

			case "chargeback":

			default:
				return fmt.Errorf("unrecognized 3-item transaction type %q at line %d", typ, lineNumber)
			}
			continue
		}

		amtStr := strings.TrimSpace(row[3])
		amt, err := strconv.ParseFloat(amtStr, 64)
		if err != nil {
			return fmt.Errorf("invalid amount %q at line %d", amtStr, lineNumber)
		}
		switch typ {
		case "deposit":
			cli.available += amt
		case "withdrawal":
			cli.available -= amt
		default:
			return fmt.Errorf("unrecognized transaction type %q at line %d", typ, lineNumber)
		}
		e.transactions[strings.TrimSpace(row[2])] = transaction{typ: typ, amount: amt}
	}

	fmt.Fprintln(w, "client, available, held, total, locked")
	for id, cli := range e.clients {
		total := cli.available + cli.held
		fmt.Fprintf(w, "%s, %.4f, %.4f, %.4f, %v\n", id, cli.available, cli.held, total, cli.locked)
	}

	return nil
}
