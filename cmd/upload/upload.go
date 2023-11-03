package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

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
	dbConfig := db.AddFlags()
	var (
		formatString = flag.String("format", "dkb", "what format the csv file is in")
		filePath     = flag.String("file", "", "path to the csv file")
		dryFilePath  = flag.String(
			"dry-file",
			"",
			"dry run the script and save the transactions and holders that would be created to the json file")
	)
	flag.Parse()
	slog.Info("flags", slog.Group("flags",
		slog.String("format", *formatString),
		slog.String("file", *filePath),
		slog.String("dry-file", *dryFilePath),
	))

	if err := logConfig.InitSlogDefault(); err != nil {
		return fmt.Errorf("init slog: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

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

	if *dryFilePath != "" {
		// this is a dry run, so we just print the transactions
		// and holders that would be created
		// TODO: actually check which holders and transactions already exist
		holderMap := make(map[db.HolderIdentifier]db.CreateHolder)
		transactions := make([]db.CreateTransaction, 0)
		for _, creator := range creators {
			transaction := creator.Transaction()
			transactions = append(transactions, transaction)
			fromHolder := creator.FromHolder()
			toHolder := creator.ToHolder()
			holderMap[fromHolder.HolderIdentifier] = fromHolder
			holderMap[toHolder.HolderIdentifier] = toHolder
		}
		holders := make([]db.CreateHolder, 0)
		for _, holder := range holderMap {
			holders = append(holders, holder)
		}
		jsonData := struct {
			Holders      []db.CreateHolder      `json:"holders"`
			Transactions []db.CreateTransaction `json:"transactions"`
		}{
			Holders:      holders,
			Transactions: transactions,
		}
		jsonBytes, err := json.Marshal(jsonData)
		if err != nil {
			return fmt.Errorf("marshal json: %w", err)
		}
		if err := os.WriteFile(*dryFilePath, jsonBytes, 0644); err != nil {
			return fmt.Errorf("write json file: %w", err)
		}
	} else {
		// this is not a dry run, so we create the transactions
		// and holders
		dbConn, err := dbConfig.Open()
		if err != nil {
			return fmt.Errorf("open connection: %w", err)
		}
		defer dbConn.Close()

		// Check the connection
		err = dbConn.Ping()
		if err != nil {
			return fmt.Errorf("ping db: %w", err)
		}

		slog.Info("successfully connected")

		for _, creator := range creators {
			transaction := creator.Transaction()
			fromHolder, err := db.GetHolderByIdentifier(ctx, dbConn, transaction.FromIdentifier)
			if err != nil {
				fromHolder, err = db.InsertHolder(ctx, dbConn, creator.FromHolder())
				if err != nil {
					return fmt.Errorf("insert from holder: %w", err)
				}
			}
			toHolder, err := db.GetHolderByIdentifier(ctx, dbConn, transaction.ToIdentifier)
			if err != nil {
				toHolder, err = db.InsertHolder(ctx, dbConn, creator.ToHolder())
				if err != nil {
					return fmt.Errorf("insert to holder: %w", err)
				}
			}
			slog.Debug("inserting transaction", slog.Group("transaction",
				slog.String("from", fromHolder.Name),
				slog.String("to", toHolder.Name)))
		}
	}
	return nil
}
