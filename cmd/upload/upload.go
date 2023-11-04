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
	"github.com/jmoiron/sqlx"
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

	// connect to db
	dbConn, err := dbConfig.OpenPingedConnection()
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer dbConn.Close()

	dryRunResult, err := upload.DryRun(ctx, dbConn, creators)
	if err != nil {
		return fmt.Errorf("dry run: %w", err)
	}
	slog.Debug("dry run result", slog.Group("result",
		slog.Int("existingHolders", len(dryRunResult.ExistingHolders)),
		slog.Int("holdersToCreate", len(dryRunResult.HoldersToCreate)),
		slog.Int("transactions", len(dryRunResult.Transactions)),
	))

	if *dryFilePath != "" {
		// this is a dry run, so we just save the result to the json file
		jsonBytes, err := json.Marshal(dryRunResult)
		if err != nil {
			return fmt.Errorf("json marshal dry run result: %w", err)
		}
		if err := os.WriteFile(*dryFilePath, jsonBytes, 0644); err != nil {
			return fmt.Errorf("write json file: %w", err)
		}
		return nil
	}

	return nonDryRun(ctx, dbConn, dryRunResult)
}

func nonDryRun(ctx context.Context,
	dbConn sqlx.ExtContext,
	result *upload.DryRunResult) error {

	err := result.InsertHolders(ctx, dbConn)
	if err != nil {
		return fmt.Errorf("insert holders: %w", err)
	}

	for _, transaction := range result.Transactions {
		fromHolder, ok := result.ExistingHolders[transaction.FromIdentifier]
		if !ok {
			return fmt.Errorf("from holder not found")
		}
		toHolder, ok := result.ExistingHolders[transaction.ToIdentifier]
		if !ok {
			return fmt.Errorf("to holder not found")
		}

		create := db.CreateTransaction{
			BaseTransaction: transaction.Transaction,
			FromHolderID:    fromHolder.ID,
			ToHolderID:      toHolder.ID,
		}

		_, err = db.InsertTransaction(ctx, dbConn, create)
		if err != nil {
			return fmt.Errorf("insert transaction: %w", err)
		}
	}
	return nil
}
