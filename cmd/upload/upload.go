package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Opsi/sparschwein/db"
	"github.com/Opsi/sparschwein/upload"
	"github.com/Opsi/sparschwein/upload/dkb"
	"github.com/Opsi/sparschwein/util"
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
	logConfig := util.AddLogFlags()
	_ = db.AddFlags()
	var (
		formatString = flag.String("format", "dkb", "what format the csv file is in")
		filePath     = flag.String("file", "", "path to the csv file")
		dryRun       = flag.Bool("dry", false, "dry run the script without writing to db")
	)
	flag.Parse()

	if err := logConfig.InitSlogDefault(); err != nil {
		return fmt.Errorf("init slog: %w", err)
	}

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
		// and holders that would be created
		for _, creator := range creators {
			fmt.Printf("Transaction: %v\n", creator.Transaction())
			fmt.Printf("From Holder: %v\n", creator.FromHolder())
			fmt.Printf("To Holder: %v\n", creator.ToHolder())
		}
	} else {
		// this is not a dry run, so we create the transactions
		// and holders
		return fmt.Errorf("not implemented")
	}
	return nil
}
