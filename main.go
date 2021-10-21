package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ijt/tx/engine"
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: tx transactions.csv\n")
		os.Exit(1)
	}
	f, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "tx: opening transactions file: %v\n", err)
		os.Exit(1)
	}
	e := engine.New()
	if err := e.Run(f, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "tx: %v\n", err)
		os.Exit(1)
	}
}
