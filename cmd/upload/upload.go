package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Opsi/sparschwein/upload"
	"github.com/Opsi/sparschwein/upload/dkb"
	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("load .env: %w", err)
	}

	// flags
	// TODO: add flag/env set for the postgres connection
	var (
		formatString = flag.String("format", "dkb", "what format to use")
		filePath     = flag.String("file", "", "path to file")
		dryRun       = flag.Bool("dry", false, "dry run the script without writing to db")
	)
	flag.Parse()

	// validate flags
	if *filePath == "" {
		return fmt.Errorf("file path is required")
	}

	// read file
	csvFile, err := os.ReadFile(*filePath)
	if err != nil {
		return fmt.Errorf("read csv file: %w", err)
	}

	var creators []upload.TransactionCreator
	switch *formatString {
	case "dkb":
		creators, err = dkb.ParseCSV(csvFile)
	default:
		return fmt.Errorf("unknown format: %s", *formatString)
	}

	if err != nil {
		return fmt.Errorf("parse csv: %w", err)
	}

	if *dryRun {
		// this is a dry run, so we just print the transactions
		// and ledgers that would be created
		for _, creator := range creators {
			fmt.Printf("Transaction: %v\n", creator.Transaction())
			fmt.Printf("From Ledger: %v\n", creator.FromLedger())
			fmt.Printf("To Ledger: %v\n", creator.ToLedger())
		}
	} else {
		// this is not a dry run, so we create the transactions
		// and ledgers
		return fmt.Errorf("not implemented")
	}
	return nil
}
